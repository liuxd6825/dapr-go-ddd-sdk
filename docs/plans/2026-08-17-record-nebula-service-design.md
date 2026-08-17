# record_nebula_service 设计文档

**日期**: 2026-08-17
**作者**: opencode (brainstorming)
**状态**: 已批准 — 待进入实施阶段

---

## 1. 背景与目标

现有 `RecordParquetService`（位于 `app/domain/graph/master/service/record_service.go`）将 Parquet 格式的银行流水从 S3/MinIO 导入 **Neo4j**，借助 Neo4j 的 `apoc.load.parquet + apoc.periodic.iterate` 在服务端完成批量解析与写入。

需求：为同一业务（银行流水导入）增加 **NebulaGraph V3.8.0** 版本。NebulaGraph 没有类似 `apoc.load.parquet` 的服务端 Parquet 直读能力，也无 MERGE/ON MATCH 语义，因此实现方式与 Neo4j 版有本质差异。

**目标**：新增 `RecordNebulaService`（及其 DAO / REST 层），以**流式**方式将 Parquet 从 MinIO 直接读入 Go 进程，按批通过 `vesoft nebula-go/v3` 写入 NebulaGraph V3.8.0，并保留与 Neo4j 版等价的幂等累加语义。

---

## 2. 关键约束

| 约束 | 说明 |
|---|---|
| NebulaGraph 版本 | V3.8.0（nGQL，无 Cypher）|
| Nebula 客户端 | `github.com/vesoft-inc/nebula-go/v3`（Thrift，GraphD 默认 9669）|
| Parquet 客户端 | `github.com/xitongsys/parquet-go`（支持 `io.Reader` 流式读）|
| MinIO 客户端 | 复用现有 `env.Minio.Client`（minio-go/v7）|
| 不下载到本地 | 流式直读（`minio.GetObject` 返回 `io.Reader`，`parquet-go` 直接消费）|
| 幂等性 | 用 `INSERT ... ON DUPLICATE KEY UPDATE` 实现累加 |
| 预创建 schema | Service 启动时调用 `initSchema`（CREATE ... IF NOT EXISTS）|

---

## 3. 架构概览

镜像现有 Neo4j 版的双层架构，**并行**新增 NebulaGraph 版本：

```
HTTP POST /case/{caseId}/record/nebula/import
  → RecordNebulaAPI.Import
    → RecordNebulaService.ImportFromS3(ctx, s3Path, caseId, opts)
          ├─ 校验：s3Path / caseId / tenantId（从 appctx）
          ├─ resolveS3Path()  // 返回 (bucket, key, *env.Minio)
          └─ RecordNebulaDao.ImportFromS3(ctx, params)
                ├─ initSchema()       // CREATE SPACE / TAG / EDGE IF NOT EXISTS
                ├─ minio.Client.GetObject(ctx, bucket, key, opts)
                │     → *minio.Object (implements io.Reader)
                ├─ parquet.NewGenericReader[Row](obj)  // 流式 reader
                ├─ 按 BatchSize=500 迭代 Read(buf)
                ├─ 对每批：批量 INSERT VERTEX + INSERT EDGE
                │           ON DUPLICATE KEY UPDATE
                └─ 汇总 ImportResult (rows / vertices / edges / errors)
```

---

## 4. 文件变更清单

| 路径 | 类型 | 说明 |
|---|---|---|
| `pkg/env/nebula.go` | 新增 | `Nebula` 类型 + `InitDBNebula` + `GetNebulaByKey` + `AddNebula` |
| `pkg/env/env.go` | 修改 | `Env` struct 加 `Nebula map[string]*Nebula` 字段；`Env.Init()` 调用 `InitDBNebula(env)` |
| `app/domain/graph/master/dao/record_nebula_dao.go` | 新增 | DAO 层：单例、`ImportParams`、`ImportFromS3`、`initSchema`、`insertBatch`、`BuildInsertNGQLs` |
| `app/domain/graph/master/service/record_nebula_service.go` | 新增 | Service 层：`RecordNebulaService`、`ImportOptions`、`ImportFromS3`、`resolveS3Path` |
| `app/domain/graph/master/restapi/record_nebula_api.go` | 新增 | REST 路由 `POST /case/{caseId}/record/nebula/import` |
| `app/domain/graph/master/restapi/register.go` | 修改 | 新增 `RegisterRecordNebulaApi` |
| `app/domain/graph/master/service/record_nebula_service_test.go` | 新增 | 单测 + 集成测试 |
| `go.mod` | 修改 | 新增 `vesoft-inc/nebula-go/v3`、`xitongsys/parquet-go` |

---

## 5. 组件设计

### 5.1 `pkg/env/nebula.go`（新增）

```go
package env

type Nebula struct {
    Name     string   `yaml:"-"`
    Addrs    []string `yaml:"addrs"`    // e.g. ["127.0.0.1:9669"]
    User     string   `yaml:"user"`
    Password string   `yaml:"password"`
    Space    string   `yaml:"space"`    // 默认 space 名
    PoolSize int      `yaml:"poolSize"` // 默认 10
    pool     *nebula.ConnectionPool    `yaml:"-"`
}

func InitDBNebula(env *Env)            // 初始化所有 nebula 连接池
func (env *Env) GetNebulaByKey(key string) (*Nebula, bool)
func (env *Env) AddNebula(n *Nebula)
```

**关键点**
- NebulaGraph V3.8.0 服务端地址是 **GraphD**（默认 `9669`），不是 Meta（9559）或 Storage（9779）。
- `Addrs` 是 **GraphD 地址列表**，`nebula-go` 的 `ConnectionPool` 会按一致性哈希负载。
- 复用 `env.Env` 单例模式，与 `Neo4j`/`Minio` 一致。

### 5.2 `pkg/env/env.go`（修改）

```go
type Env struct {
    // ... 现有字段 ...
    Nebula map[string]*Nebula `yaml:"nebula" json:"nebula"`
}

func (env *Env) Init() {
    // ... 现有初始化 ...
    InitDBNebula(env)
}
```

### 5.3 `dao/record_nebula_dao.go`（新增）

```go
package dao

import (
    "context"
    "fmt"
    "sync"
    "time"

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

// Row 内部流式读取时的行结构
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

func (d *RecordNebulaDao) ImportFromS3(ctx context.Context, p ImportParams) (res *ImportResult, err error) {
    defer func() { err = errors.GetRecoverError(err, recover()) }()

    // 校验 + 默认值
    // ... (同 Neo4j 版)

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
        return nil, err
    }

    // 流式：minio.GetObject → parquet reader
    obj, err := m.Client.GetObject(ctx, p.Bucket, p.Key, minio.GetObjectOptions{})
    if err != nil { return nil, err }
    defer obj.Close()

    pf, err := local.NewLocalFileReaderFromIOReader(obj, p.Key, int64(p.FileSize))
    // 或使用 parquet-go 的 NewReader 直接接受 io.Reader：
    // pr, err := reader.NewParquetReader(reader.NewParquetReader(...))
    pr, err := reader.NewParquetReader(pf, p.Key, 1)
    if err != nil { return nil, err }
    defer pr.ReadStop()

    num := int(pr.GetNumRows())
    res = &ImportResult{Total: num}
    for start := 0; start < num; start += p.BatchSize {
        end := start + p.BatchSize
        if end > num { end = num }
        batch := make([]Row, end-start)
        if err := pr.Read(&batch); err != nil { ... }

        // 重试
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
    return res, nil
}

func (d *RecordNebulaDao) initSchema(ctx, session, space string) error {
    stmts := []string{
        fmt.Sprintf("CREATE SPACE IF NOT EXISTS %s", space),
        fmt.Sprintf("USE %s", space),
        `CREATE TAG IF NOT EXISTS human(id string, name string, tenant_id string, case_id string)`,
        `CREATE TAG IF NOT EXISTS account(id string, name string, tenant_id string, case_id string)`,
        `CREATE EDGE IF NOT EXISTS owner()`,
        `CREATE EDGE IF NOT EXISTS record(
            id string, tenant_id string, case_id string,
            date timestamp, amount double, payout double,
            income double, txn_count int, ccy string
        )`,
        `CREATE TAG INDEX IF NOT EXISTS human_id_idx ON human(id)`,
        `CREATE TAG INDEX IF NOT EXISTS account_id_idx ON account(id)`,
        `CREATE EDGE INDEX IF NOT EXISTS record_id_idx ON record(id)`,
    }
    for _, s := range stmts {
        if err := session.Execute(s); err != nil {
            // IF NOT EXISTS 已存在的报错应忽略
            if !strings.Contains(err.Error(), "already exists") {
                return fmt.Errorf("init schema [%s] failed: %w", s, err)
            }
        }
    }
    return nil
}

func (d *RecordNebulaDao) insertBatch(ctx, session, p ImportParams, batch []Row) error {
    ngqls := BuildInsertNGQLs(p, batch)
    for _, q := range ngqls {
        if err := session.Execute(q); err != nil {
            return err
        }
    }
    return nil
}

// BuildInsertNGQLs 生成 INSERT VERTEX + INSERT EDGE nGQL
func BuildInsertNGQLs(p ImportParams, batch []Row) []string {
    // 1) 去重：human / account 节点 ID
    humans := make(map[string]Row)
    accounts := make(map[string]Row)
    for _, r := range batch {
        hID := r.Name
        if hID == "" { hID = r.Acct }
        nID := r.OppName
        if nID == "" { nID = r.OppAcct }
        humans[hID] = r
        humans[nID] = r
        accounts[r.Acct] = r
        accounts[r.OppAcct] = r
    }

    out := []string{}

    // INSERT VERTEX human
    var humanVals []string
    for id, r := range humans {
        humanVals = append(humanVals,
            fmt.Sprintf(`"%s":("%s", "%s", "%s", "%s")`,
                id, escape(id), r.TenantId, r.CaseId, escape(id)))
    }
    out = append(out,
        fmt.Sprintf(`INSERT VERTEX human(id, name, tenant_id, case_id) VALUES %s`,
            strings.Join(humanVals, ", ")))

    // INSERT VERTEX account
    var acctVals []string
    for id, r := range accounts {
        acctVals = append(acctVals,
            fmt.Sprintf(`"%s":("%s", "%s", "%s", "%s")`,
                id, escape(id), r.TenantId, r.CaseId, escape(id)))
    }
    out = append(out,
        fmt.Sprintf(`INSERT VERTEX account(id, name, tenant_id, case_id) VALUES %s`,
            strings.Join(acctVals, ", ")))

    // INSERT EDGE owner (human -> account)
    var ownerVals []string
    for _, r := range batch {
        hID := r.Name; if hID == "" { hID = r.Acct }
        ownerVals = append(ownerVals,
            fmt.Sprintf(`"%s"->"%s":()`, hID, r.Acct))
        nID := r.OppName; if nID == "" { nID = r.OppAcct }
        ownerVals = append(ownerVals,
            fmt.Sprintf(`"%s"->"%s":()`, nID, r.OppAcct))
    }
    out = append(out,
        fmt.Sprintf(`INSERT EDGE owner() VALUES %s`,
            strings.Join(ownerVals, ", ")))

    // INSERT EDGE record (account -> account) ON DUPLICATE KEY UPDATE
    var recordVals []string
    for _, r := range batch {
        recID := fmt.Sprintf("rec:%s:%s:%s:%s:%s",
            r.TenantId, r.CaseId, r.Acct, r.OppAcct,
            time.UnixMilli(r.Date).Format("2006-01-02"))
        recordVals = append(recordVals,
            fmt.Sprintf(`"%s"->"%s":("%s", "%s", "%s", %d, %.6f, %.6f, %.6f, 1, "%s")`,
                r.Acct, r.OppAcct, recID, r.TenantId, r.CaseId,
                r.Date, r.Amount, r.Payout, r.Income, r.Ccy))
    }
    out = append(out,
        fmt.Sprintf(`INSERT EDGE record(id, tenant_id, case_id, date, amount, payout, income, txn_count, ccy)
            VALUES %s
            ON DUPLICATE KEY UPDATE
                amount = record.amount + $-.amount,
                payout = record.payout + $-.payout,
                income  = record.income  + $-.income,
                txn_count = record.txn_count + 1`,
            strings.Join(recordVals, ", ")))
    return out
}
```

**设计要点**
- 单例模式：`var _recordNebulaDao + sync.Once`
- `env *env.Env` 通过构造函数传入（不强制单例依赖全局，便于测试）
- 流式读：`minio.GetObject` 返回 `io.ReadCloser`，`parquet-go-source/local` 提供 `NewLocalFileReaderFromIOReader` 包装
- 幂等：`INSERT EDGE ... ON DUPLICATE KEY UPDATE`，NebulaGraph V3.8.0 支持
- 累加：使用管道语法 `record.amount + $-.amount`（`-` 占位表示当前边）

### 5.4 `service/record_nebula_service.go`（新增）

```go
package service

import (
    "context"
    "fmt"
    "strings"

    "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
    "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

const (
    defaultMinioName    = "default"
    defaultBucketConfig = "importData"
    defaultBatchSize    = 500
    defaultRetries      = 1
    defaultConcurrency  = 1
    defaultNebulaName   = "default"
)

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

// ImportOptions 与 Neo4j 版字段一致
type ImportOptions struct {
    BatchSize   int
    Parallel    bool
    Concurrency int
    Retries     int
}

func (s *RecordNebulaService) ImportFromS3(ctx context.Context, s3Path, caseId string, opts ImportOptions) (*dao.ImportResult, error) {
    if s3Path == "" { return nil, errors.New("s3Path is empty") }
    if caseId == "" { return nil, errors.New("caseId is empty") }
    tenantId := appctx.GetTenantId2(ctx)
    if tenantId == "" { return nil, errors.New("tenantId is empty in context") }

    bucket, key, err := s.resolveS3Path(s3Path)
    if err != nil { return nil, err }

    params := dao.ImportParams{
        TenantId: tenantId, CaseId: caseId, Bucket: bucket, Key: key,
        BatchSize: opts.BatchSize, Parallel: opts.Parallel,
        Concurrency: opts.Concurrency, Retries: opts.Retries,
    }
    if params.BatchSize <= 0 { params.BatchSize = defaultBatchSize }
    if params.Retries <= 0 { params.Retries = defaultRetries }
    if params.Concurrency <= 0 { params.Concurrency = defaultConcurrency }

    var res *dao.ImportResult
    gp.Try(func() error {
        r, e := s.dao.ImportFromS3(ctx, params)
        if e != nil { return e }
        res = r
        return nil
    }).Catch(func(e error) { err = e }).Finally(func() {})

    return res, err
}

// resolveS3Path 解析三种形式（同 Neo4j 版），返回 (bucket, key)
func (s *RecordNebulaService) resolveS3Path(s3Path string) (bucket, key string, err error) {
    if strings.HasPrefix(s3Path, "s3://") {
        trimmed := strings.TrimPrefix(s3Path, "s3://")
        if strings.Contains(trimmed, "@") {
            // 形式 1: s3://user:pass@host/bucket/key → 跳过凭证段
            parts := strings.SplitN(trimmed, "/", 2)
            parts = strings.SplitN(parts[1], "/", 2)
            return parts[0], parts[1], nil
        }
        // 形式 2: s3://bucket/key
        minio, ok := s.env.GetMinioByKey(defaultMinioName)
        if !ok { return "", "", errors.New("env minio.default not configured") }
        parts := strings.SplitN(trimmed, "/", 2)
        return parts[0], parts[1], nil
    }
    // 形式 3: bucket/key
    minio, ok := s.env.GetMinioByKey(defaultMinioName)
    if !ok { return "", "", errors.New("env minio.default not configured") }
    bucketName, ok := minio.Buckets[defaultBucketConfig]
    if !ok { return "", "", errors.New("env minio.default.buckets.importData not configured") }
    return bucketName, s3Path, nil
}
```

**设计要点**
- 公开 API 与 `RecordParquetService` 完全一致：便于上层调用方切换
- `resolveS3Path` 替代 Neo4j 版的 `resolveS3URL`（不拼 `s3://user:pass@host`），返回 bucket/key
- 默认值常量与 Neo4j 版对齐

### 5.5 `restapi/record_nebula_api.go`（新增）

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
    return &RecordNebulaAPI{env: env, rootPath: rootPath, service: service.NewRecordNebulaService(env)}
}

type ImportRecordNebulaPath struct {
    CaseId string `path:"caseId" required:"true"`
}

type ImportRecordNebulaRequest struct {
    S3Path      string `json:"s3Path" validate:"required"`
    BatchSize   int    `json:"batchSize,omitempty"`
    Parallel    bool   `json:"parallel,omitempty"`
    Concurrency int    `json:"concurrency,omitempty"`
    Retries     int    `json:"retries,omitempty"`
}

func (s *RecordNebulaAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
    ctl := restapi.NewController(app, s.rootPath, "master.RecordNebulaAPI", s)
    ctl.Post("/case/{caseId}/record/nebula/import", "Import")
    ctl.Handle(iris.MethodOptions, "/case/{caseId}/record/nebula/import", "Check")
    return ctl
}

func (s *RecordNebulaAPI) Import(ctx context.Context, req *ImportRecordNebulaRequest, path *ImportRecordNebulaPath) (*dao.ImportResult, error) {
    opts := service.ImportOptions{
        BatchSize: req.BatchSize, Parallel: req.Parallel,
        Concurrency: req.Concurrency, Retries: req.Retries,
    }
    return s.service.ImportFromS3(ctx, req.S3Path, path.CaseId, opts)
}

func (s *RecordNebulaAPI) Check(ctx context.Context) error {
    logs.Infofmt(ctx, "graph/master/restapi/record-nebula/import:check")
    return nil
}
```

### 5.6 `restapi/register.go`（修改）

```go
func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
    // ...
    RegisterDrawApi(app, baseUrl, env)
    RegisterRecordParquetApi(app, baseUrl, env)
    RegisterRecordNebulaApi(app, baseUrl, env) // 新增
}

func RegisterRecordNebulaApi(app *iris.Application, baseUrl string, env *env.Env) {
    restapi.RegisterController(app, NewRecordNebulaAPI(env, baseUrl))
}
```

### 5.7 `go.mod`（修改）

新增依赖：
```
github.com/vesoft-inc/nebula-go/v3 v3.6.0
github.com/xitongsys/parquet-go v1.6.0
github.com/xitongsys/parquet-go-source v0.0.0-20210211071506-0aca4a9c4e62
```

---

## 6. nGQL 语义详解

### 6.1 Schema 创建

```nGQL
CREATE SPACE IF NOT EXISTS master_record           -- partition 与 replica 后续用 ALTER 调整
USE master_record
CREATE TAG IF NOT EXISTS human(id string, name string, tenant_id string, case_id string)
CREATE TAG IF NOT EXISTS account(id string, name string, tenant_id string, case_id string)
CREATE EDGE IF NOT EXISTS owner()
CREATE EDGE IF NOT EXISTS record(
    id string, tenant_id string, case_id string,
    date timestamp, amount double, payout double,
    income double, txn_count int, ccy string)
CREATE TAG INDEX IF NOT EXISTS human_id_idx ON human(id)
CREATE TAG INDEX IF NOT EXISTS account_id_idx ON account(id)
CREATE EDGE INDEX IF NOT EXISTS record_id_idx ON record(id)
```

> **重要**：NebulaGraph 中 `date` 是 `date` 类型；为兼容 Parquet 中的 timestamp，存储时使用 `timestamp` 类型。

### 6.2 节点插入

```nGQL
INSERT VERTEX human(id, name, tenant_id, case_id)
VALUES "张三":("张三", "张三", "t001", "case1")
```

### 6.3 边插入（关键：幂等累加）

```nGQL
INSERT EDGE record(id, tenant_id, case_id, date, amount, payout, income, txn_count, ccy)
VALUES "acct1"->"acct2":(
    "rec:t001:case1:acct1:acct2:2026-08-01",
    "t001", "case1", 1722470400000,
    1000.0, 0.0, 1000.0, 1, "CNY")
ON DUPLICATE KEY UPDATE
    amount = record.amount + $-.amount,
    payout = record.payout + $-.payout,
    income  = record.income  + $-.income,
    txn_count = record.txn_count + 1
```

**关键语法**（NebulaGraph V3.8.0）：
- `-.fieldname` 表示**待插入边**的对应字段（参考 `UPDATE` 语句中的 `$^.fieldname` 表示起点节点的字段语法）
- `record.fieldname` 表示**已存在边**的对应字段（用边类型名引用）
- `ON DUPLICATE KEY UPDATE` 在 INSERT 语句中支持（V3.x 引入）

---

## 7. 错误处理 & 并发

| 场景 | 处理 |
|---|---|
| `minio.GetObject` 失败 | 返回 error，`ImportResult = nil` |
| parquet 解析失败 | 返回 error |
| Nebula 连接失败 | 返回 error（DAO 层） |
| schema 已存在 | `IF NOT EXISTS` 幂等跳过；或在 catch 中检查 `already exists` 错误码 |
| 单批 INSERT 部分失败 | NebulaGraph 是事务性的，单批要么全成要么全败（除非事务不开启） |
| 单批 INSERT 全部失败 | 重试 `Retries` 次；仍失败计入 `ErrorMessages`，继续下一批 |
| ctx 超时 | 立刻中断批次循环 |
| 并发 `Concurrency > 1` | 每 goroutine 独立 `pool.GetSession`（**V1**：先实现 Concurrency=1，预留扩展点）|

**简化决策**：V1 不实现并发导入（NebulaGraph 单 INSERT 已是批量插入，单点瓶颈在于 session）。保持实现简单，预留 `Concurrency` 参数以便后续扩展。

---

## 8. 测试策略

`app/domain/graph/master/service/record_nebula_service_test.go`：

### 8.1 离线单测（无需环境）
- **TestBuildInsertNGQLs**：构造固定 `[]Row`，断言生成的 nGQL 包含关键片段：
  - `INSERT VERTEX human`
  - `INSERT VERTEX account`
  - `INSERT EDGE owner`
  - `INSERT EDGE record`
  - `ON DUPLICATE KEY UPDATE`
  - `rec:tenant:case:acct:oppAcct:date`
  - `txn_count = record.txn_count + 1`
- **TestResolveS3Path**：三种 s3Path 形式断言解析正确
- **TestEmptyPath**：空路径/空 caseId 错误路径

### 8.2 集成测试（需 MinIO + NebulaGraph 可达）
```go
//go:build integration
```
或通过 `testing.Short()` 跳过：
- **TestImportFromS3**：
  - 第一次导入：断言 `Total == 738`，`FailedOperations == 0`
  - 第二次导入：断言 `Total == 738`，`txn_count` 累加（用 MATCH + RETURN sum(txn_count) 验证）
- 测试 helper：镜像 Neo4j 版 `xtest.InitEnv_Neo4j` 模式，加 `InitEnv_Nebula`

> **集成测试环境前提**：NebulaGraph V3.8.0 已部署且开启 GraphD（默认端口 9669）。

---

## 9. 风险与缓解

| 风险 | 缓解 |
|---|---|
| `vesoft-inc/nebula-go/v3` API 频繁变动 | 使用稳定 API：`ConnectionPool`、`Session.Execute`、`session.Release`；锁版本 `v3.6.0` |
| NebulaGraph ON DUPLICATE KEY UPDATE 语法在 v3.8.0 可能未完全稳定 | 在集成测试中专门验证该语法；如不支持则降级为 `LOOKUP + UPDATE` 两阶段 |
| Parquet `Row` 结构与实际文件列名不匹配 | 用 `parquet-go` tag 显式声明每个字段名；不匹配时 Read 报错并立即 fail-fast |
| MinIO/S3 大文件网络抖动 | 单批失败重试 `Retries` 次 |
| Nebula session 长连接断开 | `pool.GetSession` 自动重连；DAO 层 catch 时记录错误并尝试重建 session |

---

## 10. 后续扩展（YAGNI：本设计不实现）

- 自动根据 Parquet schema 推断 NebulaGraph TAG/EDGE 字段类型（V1 写死 schema）
- 并发导入（V1 用 `Concurrency=1`）
- 数据回滚（V1 不支持）
- 数据校验（V1 依赖 NebulaGraph 自身的索引约束）

---

## 11. 验收标准

- [ ] `go build ./...` 成功
- [ ] `go test ./pkg/env/... ./app/domain/graph/master/service/...` 全部通过
- [ ] 离线单测 `TestBuildInsertNGQLs` 覆盖关键 nGQL 片段
- [ ] 集成测试在 MinIO + NebulaGraph V3.8.0 可达时通过
- [ ] REST 端点 `POST /case/{caseId}/record/nebula/import` 可正常导入 738 行 Parquet
- [ ] 第二次导入验证 `txn_count` 累加（幂等性）
- [ ] 与 Neo4j 版 `RecordParquetService` 行为对齐（同样的输入应得到等价的图结构）

---

## 12. 实施步骤概要

1. 更新 `go.mod`，新增依赖
2. 创建 `pkg/env/nebula.go`
3. 修改 `pkg/env/env.go`
4. 创建 `app/domain/graph/master/dao/record_nebula_dao.go`
5. 创建 `app/domain/graph/master/service/record_nebula_service.go`
6. 创建 `app/domain/graph/master/restapi/record_nebula_api.go`
7. 修改 `app/domain/graph/master/restapi/register.go`
8. 创建 `app/domain/graph/master/service/record_nebula_service_test.go`
9. 离线单测全通过
10. 集成测试验证

---

**待进入实施**：调用 `writing-plans` skill 生成详细实施计划。