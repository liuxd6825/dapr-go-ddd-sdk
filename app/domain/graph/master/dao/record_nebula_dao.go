package dao

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	nebula "github.com/vesoft-inc/nebula-go/v3"
	miniosrc "github.com/xitongsys/parquet-go-source/minio"
	"github.com/xitongsys/parquet-go/reader"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

// RecordNebulaDao Parquet -> NebulaGraph 导入 DAO
type RecordNebulaDao struct {
	env *env.Env
}

var (
	_recordNebulaDao  *RecordNebulaDao
	_recordNebulaOnce sync.Once
)

func NewRecordNebulaDao(env *env.Env) *RecordNebulaDao {
	_recordNebulaOnce.Do(func() {
		_recordNebulaDao = &RecordNebulaDao{env: env}
	})
	return _recordNebulaDao
}

// NebulaImportParams NebulaGraph 导入参数（独立类型，避免与 Parquet ImportParams 的 S3Path 凭证格式混淆）
type NebulaImportParams struct {
	TenantId    string
	CaseId      string
	Bucket      string
	Key         string
	BatchSize   int
	Parallel    bool
	Concurrency int
	Retries     int
}

// Row parquet 流式读取时的行结构
type Row struct {
	TenantId string  `parquet:"name=tenant_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	CaseId   string  `parquet:"name=case_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Acct     string  `parquet:"name=acct, type=BYTE_ARRAY, convertedtype=UTF8"`
	OppAcct  string  `parquet:"name=opp_acct, type=BYTE_ARRAY, convertedtype=UTF8"`
	Date     int64   `parquet:"name=date, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Name     string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	OppName  string  `parquet:"name=opp_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	Amount   float64 `parquet:"name=amount, type=DOUBLE"`
	Payout   float64 `parquet:"name=payout, type=DOUBLE"`
	Income   float64 `parquet:"name=income, type=DOUBLE"`
	Ccy      string  `parquet:"name=ccy, type=BYTE_ARRAY, convertedtype=UTF8"`
}

// ImportFromS3 通过流式 Parquet 读取 + Nebula nGQL 写入
func (d *RecordNebulaDao) ImportFromS3(ctx context.Context, p NebulaImportParams) (res *ImportResult, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	if p.Bucket == "" || p.Key == "" {
		return nil, errors.New("bucket and key must be non-empty")
	}
	if p.BatchSize <= 0 {
		p.BatchSize = 500
	}
	if p.Retries <= 0 {
		p.Retries = 1
	}
	if p.Concurrency <= 0 {
		p.Concurrency = 1
	}

	n, ok := d.env.GetNebulaByKey("default")
	if !ok {
		return nil, errors.New("env nebula.default not configured")
	}

	session, err := n.GetSession(ctx)
	if err != nil {
		return nil, errors.New("get nebula session failed: %v", err)
	}
	defer session.Release()

	if err := d.initSchema(ctx, session, n.Space); err != nil {
		return nil, err
	}
	if _, err := session.Execute(fmt.Sprintf("USE %s", n.Space)); err != nil {
		return nil, errors.New("USE space failed: %v", err)
	}

	minioCfg, ok := d.env.GetMinioByKey("default")
	if !ok {
		return nil, errors.New("env minio.default not configured")
	}

	pf, err := miniosrc.NewS3FileReaderWithClient(ctx, minioCfg.Client, p.Bucket, p.Key)
	if err != nil {
		return nil, errors.New("open s3 parquet reader failed: %v", err)
	}
	defer pf.Close()

	pr, err := reader.NewParquetReader(pf, p.Key, 4)
	if err != nil {
		return nil, errors.New("new parquet reader failed: %v", err)
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	res = &ImportResult{Total: num}
	start := time.Now()

	logs.InfoMsg(ctx, "RecordNebulaDao.ImportFromS3 start",
		" bucket=", p.Bucket,
		" key=", p.Key,
		" tenantId=", p.TenantId,
		" caseId=", p.CaseId,
		" rows=", num,
	)

	gp.Try(func() error {
		for off := 0; off < num; off += p.BatchSize {
			end := off + p.BatchSize
			if end > num {
				end = num
			}
			batch := make([]Row, end-off)
			if rerr := pr.Read(&batch); rerr != nil {
				res.FailedOperations += len(batch)
				res.ErrorMessages = append(res.ErrorMessages, rerr.Error())
				continue
			}

			var lastErr error
			for attempt := 0; attempt <= p.Retries; attempt++ {
				if e := d.insertBatch(ctx, session, batch); e == nil {
					res.CommittedOperations += len(batch)
					lastErr = nil
					break
				} else {
					lastErr = e
				}
			}
			if lastErr != nil {
				res.FailedOperations += len(batch)
				res.ErrorMessages = append(res.ErrorMessages, lastErr.Error())
			}
			res.Batches++
		}
		return nil
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {})

	res.TimeTakenMs = time.Since(start).Milliseconds()
	if err != nil {
		return res, err
	}
	return res, nil
}

// initSchema 幂等创建 Space / TAG / EDGE / INDEX
func (d *RecordNebulaDao) initSchema(ctx context.Context, session *nebula.Session, space string) error {
	stmts := []string{
		fmt.Sprintf("CREATE SPACE IF NOT EXISTS %s", space),
		fmt.Sprintf("USE %s", space),
		`CREATE TAG IF NOT EXISTS human(id string, name string, tenant_id string, case_id string)`,
		`CREATE TAG IF NOT EXISTS account(id string, name string, tenant_id string, case_id string)`,
		`CREATE EDGE IF NOT EXISTS owner()`,
		`CREATE EDGE IF NOT EXISTS record(id string, tenant_id string, case_id string, date timestamp, amount double, payout double, income double, txn_count int, ccy string)`,
		`CREATE TAG INDEX IF NOT EXISTS human_id_idx ON human(id)`,
		`CREATE TAG INDEX IF NOT EXISTS account_id_idx ON account(id)`,
		`CREATE EDGE INDEX IF NOT EXISTS record_id_idx ON record(id)`,
	}
	for _, s := range stmts {
		if _, err := session.Execute(s); err != nil {
			msg := err.Error()
			if !strings.Contains(msg, "already exists") && !strings.Contains(msg, "Exist") {
				return fmt.Errorf("init schema [%s] failed: %v", s, err)
			}
		}
	}
	return nil
}

// insertBatch 对单批执行 INSERT VERTEX/EDGE
func (d *RecordNebulaDao) insertBatch(ctx context.Context, session *nebula.Session, batch []Row) error {
	ngqls := BuildInsertNGQLs(NebulaImportParams{}, batch)
	for _, q := range ngqls {
		if _, err := session.Execute(q); err != nil {
			return err
		}
	}
	return nil
}

// BuildInsertNGQLs 生成 INSERT VERTEX + INSERT EDGE nGQL 列表
func BuildInsertNGQLs(p NebulaImportParams, batch []Row) []string {
	humans := make(map[string]Row)
	accounts := make(map[string]Row)
	for _, r := range batch {
		hID := r.Name
		if hID == "" {
			hID = r.Acct
		}
		nID := r.OppName
		if nID == "" {
			nID = r.OppAcct
		}
		humans[hID] = r
		humans[nID] = r
		accounts[r.Acct] = r
		accounts[r.OppAcct] = r
	}

	var out []string

	var humanVals []string
	for id, r := range humans {
		humanVals = append(humanVals, fmt.Sprintf(`"%s":("%s", "%s", "%s", "%s")`,
			escapeNGQL(id), escapeNGQL(id), r.TenantId, r.CaseId, escapeNGQL(id)))
	}
	if len(humanVals) > 0 {
		out = append(out, fmt.Sprintf(
			`INSERT VERTEX human(id, name, tenant_id, case_id) VALUES %s`,
			strings.Join(humanVals, ", ")))
	}

	var acctVals []string
	for id, r := range accounts {
		acctVals = append(acctVals, fmt.Sprintf(`"%s":("%s", "%s", "%s", "%s")`,
			escapeNGQL(id), escapeNGQL(id), r.TenantId, r.CaseId, escapeNGQL(id)))
	}
	if len(acctVals) > 0 {
		out = append(out, fmt.Sprintf(
			`INSERT VERTEX account(id, name, tenant_id, case_id) VALUES %s`,
			strings.Join(acctVals, ", ")))
	}

	var ownerVals []string
	for _, r := range batch {
		hID := r.Name
		if hID == "" {
			hID = r.Acct
		}
		nID := r.OppName
		if nID == "" {
			nID = r.OppAcct
		}
		ownerVals = append(ownerVals, fmt.Sprintf(`"%s"->"%s":()`, escapeNGQL(hID), escapeNGQL(r.Acct)))
		ownerVals = append(ownerVals, fmt.Sprintf(`"%s"->"%s":()`, escapeNGQL(nID), escapeNGQL(r.OppAcct)))
	}
	if len(ownerVals) > 0 {
		out = append(out, fmt.Sprintf(
			`INSERT EDGE owner() VALUES %s`,
			strings.Join(ownerVals, ", ")))
	}

	var recordVals []string
	for _, r := range batch {
		txDate := time.UnixMilli(r.Date).Format("2006-01-02")
		recID := fmt.Sprintf("rec:%s:%s:%s:%s:%s",
			r.TenantId, r.CaseId, r.Acct, r.OppAcct, txDate)
		recordVals = append(recordVals, fmt.Sprintf(
			`"%s"->"%s":("%s", "%s", "%s", %d, %.6f, %.6f, %.6f, 1, "%s")`,
			escapeNGQL(r.Acct), escapeNGQL(r.OppAcct),
			recID, r.TenantId, r.CaseId,
			r.Date, r.Amount, r.Payout, r.Income, r.Ccy))
	}
	if len(recordVals) > 0 {
		out = append(out, fmt.Sprintf(
			`INSERT EDGE record(id, tenant_id, case_id, date, amount, payout, income, txn_count, ccy) VALUES %s ON DUPLICATE KEY UPDATE amount = record.amount + $-.amount, payout = record.payout + $-.payout, income = record.income + $-.income, txn_count = record.txn_count + 1`,
			strings.Join(recordVals, ", ")))
	}
	return out
}

func escapeNGQL(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}