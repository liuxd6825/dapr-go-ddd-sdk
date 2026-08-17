# record_nebula_service Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 新增 `RecordNebulaService`（及 DAO/REST 层），将 Parquet 银行流水从 MinIO 流式导入到 NebulaGraph V3.8.0，保持与 Neo4j 版相同的对外 API。

**Tech Stack:** Go 1.25 + `vesoft-inc/nebula-go/v3` + `xitongsys/parquet-go` + `minio-go/v7`

---

## Task 1: 添加依赖

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

**Step 1: 添加 nebula-go 依赖**

```bash
go get github.com/vesoft-inc/nebula-go/v3@v3.6.0
```

**Step 2: 添加 parquet-go 依赖**

```bash
go get github.com/xitongsys/parquet-go@v1.6.0
go get github.com/xitongsys/parquet-go-source@latest
```

**Step 3: 验证依赖写入 go.mod**

Run: `grep -E "vesoft|parquet" go.mod`
Expected: 看到 `github.com/vesoft-inc/nebula-go/v3` 和 `github.com/xitongsys/parquet-go` 行

**Step 4: 编译验证**

Run: `go build ./...`
Expected: BUILD OK (零退出码)

**Step 5: Commit**

```bash
git add go.mod go.sum
git commit -m "feat(nebula): add vesoft nebula-go and parquet-go dependencies"
```

---

## Task 2: 新增 `pkg/env/nebula.go`

**Files:**
- Create: `pkg/env/nebula.go`

**Step 1: 编写 Nebula 类型 + Init/Add/Get 函数**

```go
package env

import (
    "fmt"
    "strings"

    nebula "github.com/vesoft-inc/nebula-go/v3"
)

type Nebula struct {
    Name     string   `yaml:"-"`
    Addrs    []string `yaml:"addrs"`
    User     string   `yaml:"user"`
    Password string   `yaml:"password"`
    Space    string   `yaml:"space"`
    PoolSize int      `yaml:"poolSize"`
    pool     *nebula.ConnectionPool `yaml:"-"`
}

func InitDBNebula(env *Env) {
    if env.Nebula == nil {
        env.Nebula = map[string]*Nebula{}
        return
    }
    for dbKey, cfg := range env.Nebula {
        if len(cfg.Addrs) == 0 {
            continue
        }
        if cfg.PoolSize <= 0 {
            cfg.PoolSize = 10
        }
        pool, err := nebula.NewConnectionPool(cfg.Addrs, nebula.PoolConfig{
            MaxIdleConns: cfg.PoolSize,
        }, nebula.DefaultLogger{})
        if err != nil {
            panic(fmt.Sprintf("nebula 连接池创建失败: %v", err))
        }
        cfg.pool = pool
        key := strings.ToLower(dbKey)
        cfg.Name = key
    }
}

func (env *Env) GetNebulaByKey(dbKey string) (*Nebula, bool) {
    n := env.Nebula[dbKey]
    return n, n != nil
}

func (env *Env) AddNebula(cfg *Nebula) {
    if cfg == nil {
        return
    }
    env.Nebula[cfg.Name] = cfg
}
```

**Step 2: 编译验证**

Run: `go build ./pkg/env/...`
Expected: BUILD OK

**Step 3: Commit**

```bash
git add pkg/env/nebula.go
git commit -m "feat(env): add Nebula config type and connection pool"
```

---

## Task 3: 修改 `pkg/env/env.go` 注册 Nebula

**Files:**
- Modify: `pkg/env/env.go:21-37` (Env struct 加 Nebula 字段)
- Modify: `pkg/env/env.go:57-75` (NewEnv 初始化 Nebula map)
- Modify: `pkg/env/env.go:77-103` (Env.Init 加 InitDBNebula 调用)

**Step 1: 在 Env struct 中加 Nebula 字段**

查找 `Minio     map[string]*Minio    \`yaml:"minio" json:"minio"\``, 在其后插入：
```go
	Nebula     map[string]*Nebula    `yaml:"nebula" json:"nebula"`
```

**Step 2: 在 NewEnv 中初始化 Nebula map**

查找 `Minio:     map[string]*Minio{},`, 在其后插入：
```go
		Nebula:     map[string]*Nebula{},
```

**Step 3: 在 Env.Init() 中调用 InitDBNebula**

查找 `initMinio(env)`, 在其后插入：
```go
	InitDBNebula(env)
```

**Step 4: 编译验证**

Run: `go build ./pkg/env/...`
Expected: BUILD OK

**Step 5: Commit**

```bash
git add pkg/env/env.go
git commit -m "feat(env): register Nebula in Env struct and init pipeline"
```

---

## Task 4: 编写 DAO 单元测试骨架（TDD: 先写失败测试）

**Files:**
- Create: `app/domain/graph/master/dao/record_nebula_dao_test.go`

**Step 1: 写失败的 BuildInsertNGQLs 单测**

```go
package dao

import (
    "strings"
    "testing"
)

func TestRecordNebulaDao_BuildInsertNGQLs(t *testing.T) {
    batch := []Row{
        {
            TenantId: "t001",
            CaseId:   "case1",
            Acct:     "acct1",
            OppAcct:  "acct2",
            Date:     1722470400000,
            Name:     "张三",
            OppName:  "李四",
            Amount:   1000.0,
            Payout:   0.0,
            Income:   1000.0,
            Ccy:      "CNY",
        },
    }
    p := ImportParams{
        TenantId:  "t001",
        CaseId:    "case1",
        BatchSize: 500,
    }
    ngqls := BuildInsertNGQLs(p, batch)

    if len(ngqls) == 0 {
        t.Fatal("expected non-empty nGQL list")
    }
    joined := strings.Join(ngqls, "\n")
    checks := []string{
        "INSERT VERTEX human",
        "INSERT VERTEX account",
        "INSERT EDGE owner",
        "INSERT EDGE record",
        "ON DUPLICATE KEY UPDATE",
        "txn_count = record.txn_count + 1",
        `"acct1"->"acct2"`,
    }
    for _, want := range checks {
        if !strings.Contains(joined, want) {
            t.Errorf("nGQL missing fragment: %q\n--- generated ---\n%s", want, joined)
        }
    }
}
```

**Step 2: 运行测试确认失败**

Run: `go test ./app/domain/graph/master/dao/ -run TestRecordNebulaDao_BuildInsertNGQLs -v`
Expected: FAIL（因为 `BuildInsertNGQLs` 和 `Row` 类型不存在）

---

## Task 5: 实现 DAO `record_nebula_dao.go`

**Files:**
- Create: `app/domain/graph/master/dao/record_nebula_dao.go`

**Step 1: 编写完整 DAO 实现**

```go
package dao

import (
    "context"
    "fmt"
    "strings"
    "sync"
    "time"

    "github.com/minio/minio-go/v7"
    nebula "github.com/vesoft-inc/nebula-go/v3"
    "github.com/xitongsys/parquet-go-source/local"
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

// ImportParams 导入参数
type ImportParams struct {
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
func (d *RecordNebulaDao) ImportFromS3(ctx context.Context, p ImportParams) (res *ImportResult, err error) {
    defer func() {
        err = errors.GetRecoverError(err, recover())
    }()

    if p.Bucket == "" || p.Key == "" {
        return nil, errors.New("bucket/key is empty")
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

    session, err := n.pool.GetSession(ctx, n.User, n.Password)
    if err != nil {
        return nil, errors.New("get nebula session failed: %w", err)
    }
    defer session.Release()

    if err := d.initSchema(ctx, session, n.Space); err != nil {
        return nil, err
    }
    if err := session.Execute(fmt.Sprintf("USE %s", n.Space)); err != nil {
        return nil, errors.New("USE space failed: %w", err)
    }

    minioCfg, ok := d.env.GetMinioByKey("default")
    if !ok {
        return nil, errors.New("env minio.default not configured")
    }
    obj, err := minioCfg.Client.GetObject(ctx, p.Bucket, p.Key, minio.GetObjectOptions{})
    if err != nil {
        return nil, errors.New("get s3 object failed: %w", err)
    }
    defer obj.Close()

    stat, err := obj.Stat()
    if err != nil {
        return nil, errors.New("stat s3 object failed: %w", err)
    }
    fileSize := stat.Size

    pf, err := local.NewLocalFileReaderFromIOReader(obj, p.Key, fileSize)
    if err != nil {
        return nil, errors.New("wrap parquet reader failed: %w", err)
    }
    pr, err := reader.NewParquetReader(pf, p.Key, 4)
    if err != nil {
        return nil, errors.New("new parquet reader failed: %w", err)
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
            if e := d.insertBatch(ctx, session, p, batch); e == nil {
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
    res.TimeTakenMs = time.Since(start).Milliseconds()

    gp.Try(func() error { return nil }).Catch(func(e error) { err = e }).Finally(func() {})
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
        if err := session.Execute(s); err != nil {
            msg := err.Error()
            if !strings.Contains(msg, "already exists") && !strings.Contains(msg, "Exist") {
                return fmt.Errorf("init schema [%s] failed: %w", s, err)
            }
        }
    }
    return nil
}

// insertBatch 对单批执行 INSERT VERTEX/EDGE
func (d *RecordNebulaDao) insertBatch(ctx context.Context, session *nebula.Session, p ImportParams, batch []Row) error {
    ngqls := BuildInsertNGQLs(p, batch)
    for _, q := range ngqls {
        if err := session.Execute(q); err != nil {
            return err
        }
    }
    return nil
}

// BuildInsertNGQLs 生成 INSERT VERTEX + INSERT EDGE nGQL 列表
func BuildInsertNGQLs(p ImportParams, batch []Row) []string {
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
```

**Step 2: 运行 DAO 单测确认通过**

Run: `go test ./app/domain/graph/master/dao/ -run TestRecordNebulaDao_BuildInsertNGQLs -v`
Expected: PASS

**Step 3: 编译整个项目**

Run: `go build ./...`
Expected: BUILD OK

**Step 4: Commit**

```bash
git add app/domain/graph/master/dao/record_nebula_dao.go app/domain/graph/master/dao/record_nebula_dao_test.go
git commit -m "feat(nebula): add RecordNebulaDao with BuildInsertNGQLs"
```

---

## Task 6: 编写 Service 失败测试

**Files:**
- Create: `app/domain/graph/master/service/record_nebula_service_test.go`

**Step 1: 编写空路径 + 三种 s3Path 解析测试**

```go
package service

import (
    "context"
    "testing"

    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func TestRecordNebulaService_EmptyPath(t *testing.T) {
    svc := &RecordNebulaService{dao: nil, env: &env.Env{}}
    ctx := context.Background()
    if _, err := svc.ImportFromS3(ctx, "", "case1", ImportOptions{}); err == nil {
        t.Fatal("expected error for empty s3Path")
    }
    if _, err := svc.ImportFromS3(ctx, "s3://x", "", ImportOptions{}); err == nil {
        t.Fatal("expected error for empty caseId")
    }
}

func TestRecordNebulaService_ResolveS3Path(t *testing.T) {
    e := &env.Env{Minio: map[string]*env.Minio{
        "default": {
            AccessKey: "ak", SecretKey: "sk", Endpoint: "host:9000",
            Buckets: map[string]string{"importData": "import-data"},
        },
    }}
    svc := &RecordNebulaService{env: e}

    cases := []struct {
        in       string
        wantBk   string
        wantKey  string
        wantErr  bool
    }{
        {"s3://ak:sk@host/bucket/key1.parquet", "bucket", "key1.parquet", false},
        {"s3://bucket/key2.parquet", "bucket", "key2.parquet", false},
        {"key3.parquet", "import-data", "key3.parquet", false},
        {"", "", "", true},
    }
    for _, c := range cases {
        bk, k, err := svc.resolveS3Path(c.in)
        if c.wantErr {
            if err == nil { t.Errorf("input %q: expected error", c.in) }
            continue
        }
        if err != nil { t.Fatalf("input %q: %v", c.in, err) }
        if bk != c.wantBk { t.Errorf("bucket: got %q want %q", bk, c.wantBk) }
        if k != c.wantKey { t.Errorf("key: got %q want %q", k, c.wantKey) }
    }
}
```

**Step 2: 运行测试确认失败**

Run: `go test ./app/domain/graph/master/service/ -run TestRecordNebulaService_EmptyPath -v`
Expected: FAIL（`RecordNebulaService` 类型未定义）

---

## Task 7: 实现 Service `record_nebula_service.go`

**Files:**
- Create: `app/domain/graph/master/service/record_nebula_service.go`

**Step 1: 编写 Service**

```go
package service

import (
    "context"
    "strings"

    "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

const (
    nebulaMinioName    = "default"
    nebulaBucketConfig = "importData"
    nebulaBatchSize    = 500
    nebulaRetries      = 1
    nebulaConcurrency  = 1
)

// RecordNebulaService 银行流水 parquet -> NebulaGraph 导入业务编排
type RecordNebulaService struct {
    dao *dao.RecordNebulaDao
    env *env.Env
}

func NewRecordNebulaService(env *env.Env) *RecordNebulaService {
    return &RecordNebulaService{
        dao: dao.NewRecordNebulaDao(env),
        env: env,
    }
}

// ImportFromS3 导入指定 s3Path 到 NebulaGraph
func (s *RecordNebulaService) ImportFromS3(ctx context.Context, s3Path, caseId string, opts ImportOptions) (*dao.ImportResult, error) {
    if s3Path == "" {
        return nil, errors.New("s3Path is empty")
    }
    if caseId == "" {
        return nil, errors.New("caseId is empty")
    }
    tenantId := appctx.GetTenantId2(ctx)
    if tenantId == "" {
        return nil, errors.New("tenantId is empty in context")
    }

    bucket, key, err := s.resolveS3Path(s3Path)
    if err != nil {
        return nil, err
    }

    params := dao.ImportParams{
        TenantId:    tenantId,
        CaseId:      caseId,
        Bucket:      bucket,
        Key:         key,
        BatchSize:   opts.BatchSize,
        Parallel:    opts.Parallel,
        Concurrency: opts.Concurrency,
        Retries:     opts Retries,
    }
    if params.BatchSize <= 0 {
        params.BatchSize = nebulaBatchSize
    }
    if params.Retries <= 0 {
        params.Retries = nebulaRetries
    }
    if params.Concurrency <= 0 {
        params.Concurrency = nebulaConcurrency
    }

    var res *dao.ImportResult
    gp.Try(func() error {
        r, e := s.dao.ImportFromS3(ctx, params)
        if e != nil {
            return e
        }
        res = r
        return nil
    }).Catch(func(e error) {
        err = e
    }).Finally(func() {})

    return res, err
}

// resolveS3Path 解析三种形式返回 (bucket, key)
func (s *RecordNebulaService) resolveS3Path(s3Path string) (bucket, key string, err error) {
    if strings.HasPrefix(s3Path, "s3://") {
        trimmed := strings.TrimPrefix(s3Path, "s3://")
        if strings.Contains(trimmed, "@") {
            // 形式 1: s3://user:pass@host/bucket/key
            at := strings.Index(trimmed, "@")
            rest := trimmed[at+1:]
            parts := strings.SplitN(rest, "/", 2)
            if len(parts) != 2 {
                return "", "", errors.New("s3 path must include host/bucket/key, got: %s", s3Path)
            }
            return parts[0], parts[1], nil
        }
        // 形式 2: s3://bucket/key
        minioCfg, ok := s.env.GetMinioByKey(nebulaMinioName)
        if !ok {
            return "", "", errors.New("env minio.default not configured")
        }
        _ = minioCfg
        parts := strings.SplitN(trimmed, "/", 2)
        if len(parts) != 2 {
            return "", "", errors.New("s3 path must include bucket/key, got: %s", s3Path)
        }
        return parts[0], parts[1], nil
    }

    // 形式 3: bucket/key
    minioCfg, ok := s.env.GetMinioByKey(nebulaMinioName)
    if !ok {
        return "", "", errors.New("env minio.default not configured")
    }
    bucketName, ok := minioCfg.Buckets[nebulaBucketConfig]
    if !ok {
        return "", "", errors.New("env minio.default.buckets.importData not configured")
    }
    return bucketName, s3Path, nil
}
```

**Step 2: 运行 Service 单测**

Run: `go test ./app/domain/graph/master/service/ -run TestRecordNebulaService -v`
Expected: PASS

**Step 3: 编译验证**

Run: `go build ./...`
Expected: BUILD OK

**Step 4: Commit**

```bash
git add app/domain/graph/master/service/record_nebula_service.go app/domain/graph/master/service/record_nebula_service_test.go
git commit -m "feat(nebula): add RecordNebulaService with s3 path resolution"
```

---

## Task 8: 新增 REST API `record_nebula_api.go`

**Files:**
- Create: `app/domain/graph/master/restapi/record_nebula_api.go`

**Step 1: 编写 REST API**

```go
package restapi

import (
    "context"

    "github.com/kataras/iris/v12"
    "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
    "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordNebulaAPI struct {
    env      *env.Env
    rootPath string
    service  *service.RecordNebulaService
}

func NewRecordNebulaAPI(env *env.Env, rootPath string) *RecordNebulaAPI {
    return &RecordNebulaAPI{
        env:      env,
        rootPath: rootPath,
        service:  service.NewRecordNebulaService(env),
    }
}

type ImportRecordNebulaPath struct {
    CaseId string `path:"caseId" required:"true"`
}

type ImportRecordNebulaRequest struct {
    S3Path      string `json:"s3Path" validate:"required" desc:"s3://bucket/key 或 bucket/key"`
    BatchSize   int    `json:"batchSize,omitempty" desc:"批次大小, 默认 500"`
    Parallel    bool   `json:"parallel,omitempty" desc:"是否并行"`
    Concurrency int    `json:"concurrency,omitempty" desc:"并发数, 默认 1"`
    Retries     int    `json:"retries,omitempty" desc:"重试次数, 默认 1"`
}

func (s *RecordNebulaAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
    ctl := restapi.NewController(app, s.rootPath, "master.RecordNebulaAPI", s)
    ctl.Post("/case/{caseId}/record/nebula/import", "Import")
    ctl.Handle(iris.MethodOptions, "/case/{caseId}/record/nebula/import", "Check")
    return ctl
}

func (s *RecordNebulaAPI) Import(ctx context.Context, req *ImportRecordNebulaRequest, path *ImportRecordNebulaPath) (*dao.ImportResult, error) {
    opts := service.ImportOptions{
        BatchSize:   req.BatchSize,
        Parallel:    req.Parallel,
        Concurrency: req.Concurrency,
        Retries:     req.Retries,
    }
    return s.service.ImportFromS3(ctx, req.S3Path, path.CaseId, opts)
}

func (s *RecordNebulaAPI) Check(ctx context.Context) error {
    logs.Infofmt(ctx, "graph/master/restapi/record-nebula/import:check")
    return nil
}
```

**Step 2: 编译验证**

Run: `go build ./app/domain/graph/master/restapi/...`
Expected: BUILD OK

---

## Task 9: 注册 REST API

**Files:**
- Modify: `app/domain/graph/master/restapi/register.go`

**Step 1: 在 RegisterAllApi 中加调用**

查找 `RegisterRecordParquetApi(app, baseUrl, env)`, 在其后插入：
```go
	RegisterRecordNebulaApi(app, baseUrl, env)
```

**Step 2: 加 RegisterRecordNebulaApi 函数**

查找 `RegisterRecordParquetApi` 函数定义, 在其后追加：
```go
func RegisterRecordNebulaApi(app *iris.Application, baseUrl string, env *env.Env) {
    restapi.RegisterController(app, NewRecordNebulaAPI(env, baseUrl))
}
```

**Step 3: 编译验证**

Run: `go build ./...`
Expected: BUILD OK

**Step 4: Commit**

```bash
git add app/domain/graph/master/restapi/record_nebula_api.go app/domain/graph/master/restapi/register.go
git commit -m "feat(nebula): register record/nebula/import REST endpoint"
```

---

## Task 10: 端到端单测

**Step 1: 跑全部新增测试**

Run: `go test ./pkg/env/... ./app/domain/graph/master/dao/ ./app/domain/graph/master/service/ -run "Nebula|RecordNebula" -v`
Expected: PASS

**Step 2: 跑 build**

Run: `go build ./...`
Expected: BUILD OK

---

## Task 11: 集成测试（需 MinIO + NebulaGraph 可达）

**Files:**
- Create: `app/domain/graph/master/service/record_nebula_service_test.go` (在 Task 6 已建, 现追加)

**Step 1: 在测试文件末尾追加集成测试**

```go
func TestRecordNebulaService_ImportFromS3(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in -short mode")
    }

    envVal := xtest2.InitEnv_Nebula(xtest2.NebulaRemoveOption)
    ctx := xtest2.NewContext()
    svc := NewRecordNebulaService(envVal)

    s3Path := "s3://" + testS3Access + ":" + testS3Secret + "@" +
        testS3Host + "/" + testBucket + "/" + testParquetKey

    res, err := svc.ImportFromS3(ctx, s3Path, testCaseId, ImportOptions{BatchSize: 500})
    if err != nil {
        t.Fatalf("first ImportFromS3 failed: %v", err)
    }
    if res.Total != 738 {
        t.Errorf("expected total=738, got %d", res.Total)
    }
    if res.FailedOperations != 0 {
        t.Errorf("expected failed=0, got %d (errors: %v)", res.FailedOperations, res.ErrorMessages)
    }
}
```

**Step 2: 在 `pkg/xtest/` 中实现 InitEnv_Nebula helper**

按现有 `InitEnv_Neo4j` 模式镜像实现（`pkg/xtest/nebula_db.go`）：
```go
package xtest

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"

const (
    NebulaDBKey      = "default"
    NebulaHostLocal  = "localhost"
    NebulaHostRemote = "192.168.120.224"
)

func GetNebulaRemoteCfg() *env.Nebula {
    return &env.Nebula{
        Name:     NebulaDBKey,
        Addrs:    []string{NebulaHostRemote + ":9669"},
        User:     "root",
        Password: "nebula",
        Space:    "test_default",
        PoolSize: 10,
    }
}
```

并在 `xtest.InitEnv` 中追加 `envVal.AddNebula(GetNebulaRemoteCfg())`。

**Step 3: 运行集成测试**

Run: `go test ./app/domain/graph/master/service/ -run TestRecordNebulaService_ImportFromS3 -v`
Expected: PASS（前提：MinIO + NebulaGraph V3.8.0 可达）

**Step 4: Commit**

```bash
git add app/domain/graph/master/service/record_nebula_service_test.go pkg/xtest/nebula_db.go pkg/xtest/init.go
git commit -m "test(nebula): add integration test and xtest helper"
```

---

## 完成验证清单

- [ ] `go build ./...` 成功
- [ ] `go test ./pkg/env/... ./app/domain/graph/master/dao/ ./app/domain/graph/master/service/ -run "Nebula" -v` 全部通过
- [ ] 设计文档保留：`docs/plans/2026-08-17-record-nebula-service-design.md`
- [ ] 集成测试在 MinIO + NebulaGraph V3.8.0 可达时通过
- [ ] REST 端点 `POST /case/{caseId}/record/nebula/import` 注册成功
- [ ] 与 Neo4j 版公开 API 对齐：`ImportFromS3(ctx, s3Path, caseId, ImportOptions) (*ImportResult, error)`

## 提交记录总览

1. `feat(nebula): add vesoft nebula-go and parquet-go dependencies`
2. `feat(env): add Nebula config type and connection pool`
3. `feat(env): register Nebula in Env struct and init pipeline`
4. `feat(nebula): add RecordNebulaDao with BuildInsertNGQLs`
5. `feat(nebula): add RecordNebulaService with s3 path resolution`
6. `feat(nebula): register record/nebula/import REST endpoint`
7. `test(nebula): add integration test and xtest helper`