package dao

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

// RecordParquetDao
// @Description: parquet -> Neo4j 导入 DAO (apoc.periodic.iterate + apoc.load.parquet)
type RecordParquetDao struct {
	*Base[*model.Record]
}

var (
	_recordParquetDao  *RecordParquetDao
	_recordParquetOnce sync.Once
)

// NewRecordParquetDao 单例创建
func NewRecordParquetDao() *RecordParquetDao {
	_recordParquetOnce.Do(func() {
		dbSch := dbschema.NewDBSchemaWithStruct("record", &model.Record{}, "record")
		nodeCfg := &dao.DaoConfig{
			DBKey:              config.Neo4jDBKey,
			GraphType:          idao.GraphType_Rel,
			IsCancelModified:   true,
			IsCancelSoftDelete: true,
		}
		newDao := dao.NewDao[*model.Record](nodeCfg)
		_recordParquetDao = &RecordParquetDao{
			Base: &Base[*model.Record]{DBSchema: dbSch, Dao: newDao},
		}
	})
	return _recordParquetDao
}

// ImportParams 导入参数
type ImportParams struct {
	TenantId    string // 租户 ID
	CaseId      string // 案件 ID
	S3Path      string // 完整 s3:// 路径 (含凭证)
	BatchSize   int    // apoc.periodic.iterate 批大小 (默认 500)
	Parallel    bool   // 是否并行 (默认 false)
	Concurrency int    // 并发数 (默认 1)
	Retries     int    // 重试次数 (默认 1)
}

// ImportResult apoc.periodic.iterate YIELD 解析结果
type ImportResult struct {
	Batches             int      `json:"batches"`
	Total               int      `json:"total"`
	CommittedOperations int      `json:"committedOperations"`
	FailedOperations    int      `json:"failedOperations"`
	TimeTakenMs         int64    `json:"timeTakenMs"`
	ErrorMessages       []string `json:"errorMessages,omitempty"`
}

// ImportFromS3 通过 apoc.load.parquet + apoc.periodic.iterate 导入
func (d *RecordParquetDao) ImportFromS3(ctx context.Context, p ImportParams) (res *ImportResult, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	if p.S3Path == "" {
		return nil, errors.New("s3Path is empty")
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

	cypher := BuildImportCypher(p)
	params := map[string]any{
		"s3Path":      p.S3Path,
		"batchSize":   p.BatchSize,
		"parallel":    p.Parallel,
		"concurrency": p.Concurrency,
		"retries":     p.Retries,
	}

	logs.InfoMsg(ctx, "RecordParquetDao.ImportFromS3 start",
		" s3Path=", p.S3Path,
		" tenantId=", p.TenantId,
		" caseId=", p.CaseId,
		" batchSize=", p.BatchSize,
	)

	gp.Try(func() error {
		result, werr := d.GetStore().Write(ctx, cypher, params)
		if werr != nil {
			return werr
		}
		if result == nil {
			return errors.New("apoc.periodic.iterate returned nil result")
		}
		res = parseImportResult(result.Data())
		if res == nil {
			return errors.New("apoc.periodic.iterate returned no rows")
		}
		return nil
	}).Catch(func(e error) {
		err = e
	}).Finally(func() {})

	return res, err
}

// parseImportResult 从 Neo4jResult 解析 apoc.periodic.iterate 的 YIELD 字段
func parseImportResult(data map[string][]any) *ImportResult {
	if data == nil {
		return nil
	}
	res := &ImportResult{}
	if v, ok := data["batches"]; ok && len(v) > 0 {
		res.Batches = toInt(v[0])
	}
	if v, ok := data["total"]; ok && len(v) > 0 {
		res.Total = toInt(v[0])
	}
	if v, ok := data["committedOperations"]; ok && len(v) > 0 {
		res.CommittedOperations = toInt(v[0])
	}
	if v, ok := data["failedOperations"]; ok && len(v) > 0 {
		res.FailedOperations = toInt(v[0])
	}
	if v, ok := data["timeTaken"]; ok && len(v) > 0 {
		res.TimeTakenMs = toInt64(v[0])
	}
	if v, ok := data["errorMessages"]; ok && len(v) > 0 {
		if list, ok := v[0].([]any); ok && len(list) > 0 {
			res.ErrorMessages = make([]string, 0, len(list))
			for _, e := range list {
				res.ErrorMessages = append(res.ErrorMessages, fmt.Sprintf("%v", e))
			}
		}
	}
	return res
}

// BuildImportCypher 生成 apoc.periodic.iterate 完整 Cypher
// 内层 Cypher 与 cypher-shell 跑通版一致:
//   - 字符串 relId 替代 md5
//   - date(datetime(row.date)) 替代 date(row.date)
//   - MERGE + ON CREATE / ON MATCH 增量累加
func BuildImportCypher(p ImportParams) string {
	parallel := "false"
	if p.Parallel {
		parallel = "true"
	}
	concurrency := p.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	inner := `
		WITH row,
			row.tenant_id  AS tenantId,
			row.case_id    AS caseId,
			row.acct       AS acct,
			row.opp_acct   AS oppAcct,
			date(datetime(row.date))        AS txDate,
			row.name                        AS name,
			row.opp_name                    AS oppName,
			coalesce(toFloat(row.amount), 0.0) AS amount,
			coalesce(toFloat(row.payout), 0.0) AS payout,
			coalesce(toFloat(row.income),  0.0) AS income,
			row.ccy AS ccy
		WHERE tenantId IS NOT NULL AND caseId IS NOT NULL
			AND acct IS NOT NULL AND oppAcct IS NOT NULL AND txDate IS NOT NULL
		WITH *, ':tenant_' + tenantId + ':case_' + caseId + ':record' AS labels
		MERGE (n1:human{id: coalesce(name, acct)})
			ON CREATE SET n1.name = name, n1.tenant_id = tenantId, n1.case_id = caseId
		MERGE (a1:account{id: acct})
			ON CREATE SET a1.name = acct, a1.tenant_id = tenantId, a1.case_id = caseId
		MERGE (n1)-[:owner]->(a1)
		MERGE (n2:human{id: coalesce(oppName, oppAcct)})
			ON CREATE SET n2.name = oppName, n2.tenant_id = tenantId, n2.case_id = caseId
		MERGE (a2:account{id: oppAcct})
			ON CREATE SET a2.name = oppAcct, a2.tenant_id = tenantId, a2.case_id = caseId
		MERGE (n2)-[:owner]->(a2)
		WITH a1, a2, tenantId, caseId, acct, oppAcct, txDate,
			amount, payout, income, ccy
		MERGE (a1)-[r:record{
			id: 'rec:' + tenantId + ':' + caseId + ':' + acct + ':' + oppAcct + ':' + toString(txDate)
		}]->(a2)
		ON CREATE SET
			r.created_time = timestamp(),
			r.tenant_id    = tenantId,
			r.case_id      = caseId,
			r.acct         = acct,
			r.opp_acct     = oppAcct,
			r.date         = txDate,
			r.amount       = amount,
			r.payout       = payout,
			r.income       = income,
			r.txn_count    = 1,
			r.ccy          = ccy
		ON MATCH SET
			r.amount       = coalesce(r.amount, 0.0) + amount,
			r.payout       = coalesce(r.payout, 0.0) + payout,
			r.income       = coalesce(r.income,  0.0) + income,
			r.txn_count    = coalesce(r.txn_count, 0) + 1,
			r.updated_time = timestamp()`

	cypher := fmt.Sprintf(`
		CALL apoc.periodic.iterate(
			"CALL apoc.load.parquet($s3Path) YIELD value AS row RETURN row",
			%q,
			{batchSize: $batchSize, parallel: %s, iterateList: true,
			 concurrency: %d, retries: $retries}
		) YIELD batches, total, committedOperations, failedOperations,
			timeTaken, errorMessages
		RETURN batches, total, committedOperations, failedOperations,
			timeTaken, errorMessages`,
		inner, parallel, concurrency,
	)
	return cypher
}

// BuildS3URL 拼接 s3://user:pass@host:port/bucket/key 格式 URL
func BuildS3URL(m *env.Minio, bucket, key string) string {
	bucket = strings.TrimLeft(bucket, "/")
	key = strings.TrimLeft(key, "/")
	return fmt.Sprintf("s3://%s:%s@%s/%s/%s",
		m.AccessKey, m.SecretKey, m.Endpoint, bucket, key)
}

// toInt / toInt64 安全转 int (apoc.periodic.iterate 返回的可能是 int64 / int / float64)
func toInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	}
	return 0
}

func toInt64(v any) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case int:
		return int64(x)
	case int32:
		return int64(x)
	case float64:
		return int64(x)
	}
	return 0
}
