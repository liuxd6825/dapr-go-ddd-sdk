# dorisloader

将结构化数据以 **Parquet 文件** 为载体、通过 Doris 的 **Stream Load** 接口批量导入到 Doris 表的工具包。

基于 **Apache parquet-go**（segmentio 维护，Apache Arrow 团队推荐），输出完全符合 Apache Parquet 规范，**Doris 原生兼容**。

## 目录

- [快速开始](#快速开始)
- [包结构](#包结构)
- [使用示例](#使用示例)
- [字段类型映射](#字段类型映射)
- [Doris 表结构要求](#doris-表结构要求)
- [关键实现细节](#关键实现细节)
- [故障排查](#故障排查)

---

## 快速开始

```go
// 1. 写出 Parquet 文件
pf, _ := dorisloader.NewParquetFile[*Record]("./record.parquet")
for _, rec := range records {
    _ = pf.Write(rec)
}
_ = pf.WriteStop() // 内部封装 flush + footer + close

// 2. 导入到 Doris（必须显式提供 Columns 与 Doris 表列顺序一致）
loader := dorisloader.NewDorisLoader(dorisloader.Config{
    BEHost:   "192.168.120.224",
    BEPort:   8040,
    Database: "master",
    Table:    "master_record1",
    Username: "admin",
    Password: "",
})
err := loader.ImportFile("./record.parquet", dorisloader.LoadOptions{
    Format:  "parquet",
    Columns: "id, tenant_id, case_id, ..., opp_name", // 与 DESC 输出列顺序一致
})
```

## 包结构

| 文件 | 作用 |
|------|------|
| `parquet_file.go` | Parquet 文件写入封装（基于 `parquet.NewGenericWriter[T]`） |
| `loader.go` | Doris Stream Load HTTP 客户端（基于 `net/http`） |
| `model.go` | `Record` 结构定义与 `*model.RecordIe` → `*Record` 转换 |
| `parquet_file_test.go` | 单元测试，覆盖 Parquet 写入 + Doris 导入全流程 |

## 使用示例

### 完整流程

```go
func importToDoris(records []*dorisloader.Record) error {
    const parquetPath = "/tmp/records.parquet"

    pf, err := dorisloader.NewParquetFile[*dorisloader.Record](parquetPath)
    if err != nil {
        return fmt.Errorf("create parquet writer: %w", err)
    }
    for _, r := range records {
        if err := pf.Write(r); err != nil {
            return fmt.Errorf("write record: %w", err)
        }
    }
    if err := pf.WriteStop(); err != nil {
        return fmt.Errorf("flush parquet: %w", err)
    }

    loader := dorisloader.NewDorisLoader(dorisloader.Config{
        BEHost:   "192.168.120.224",
        BEPort:   8040,
        Database: "master",
        Table:    "master_record1",
        Username: "admin",
        Password: "",
    })
    return loader.ImportFile(parquetPath, dorisloader.LoadOptions{
        Format:  "parquet",
        Columns: dorisTableColumns, // 与 DESC 输出列顺序一致
    })
}
```

## 字段类型映射

`Record` 字段类型与 Doris 列的对应关系（基于 `parquet-go` v0.30+ 反射推断）：

| Go 类型 | Parquet 类型 | Doris 列类型 | 备注 |
|---------|--------------|--------------|------|
| `string` | BYTE_ARRAY + UTF8 | varchar | |
| `bool` | BOOLEAN | boolean | |
| `int32` | INT32 | int | |
| `int64` | INT64 | bigint | |
| `*float64` + `optional` | DOUBLE (optional) | double | nil → NULL |
| `*time.Time` + `optional,timestamp` | INT64 + TIMESTAMP (optional) | datetime | nil → NULL |

### parquet tag 语法

极简格式（仅 `parquet:"name"` + 可选属性）：

```go
type Record struct {
    Id          string     `parquet:"id"`
    Date        *time.Time `parquet:"date,optional,timestamp"`
    Cash        bool       `parquet:"cash"`
    RowNum      int32      `parquet:"row_num"`
    Balance     *float64   `parquet:"balance,optional"`
    // ...
}
```

支持的属性：
- `optional` — 字段可为 null
- `timestamp` — 自动写为 INT64 + TIMESTAMP logical type
- `snappy` / `gzip` / `zstd` / `lz4` / `brotli` — 列级压缩覆盖

## Doris 表结构要求

### 时间相关列

`date` / `created_time` / `updated_time` 三列：
- 类型为 `DATETIME`
- 允许 NULL（默认即可；不能 NOT NULL）

### 整型列

| Doris 类型 | Go 类型 | parquet tag |
|-----------|---------|-------------|
| `int` | `int32` | `parquet:"name"` |
| `bigint` | `int64` | `parquet:"name"` |

### 数值可选列

`balance` / `payout` / `income` / `amount`：类型 `DOUBLE`，允许 NULL。

### 布尔列

`cash`：类型 `boolean`。

### 字符串列

必须为 `VARCHAR` / `STRING`。若有 `NOT NULL` 约束，调用方需保证字段非空。

## 关键实现细节

### 1. HTTP 请求不能用 chunked encoding

**Doris Stream Load 不接受 `Transfer-Encoding: chunked`**，要求 `Content-Length` header。

如果直接把 `*os.File` 作为 HTTP body，Go 会自动使用 chunked encoding，Doris 会**静默返回 0 行**（HTTP 200 + Status=Success + 0 行）。

修复：`ImportFile` 必须先用 `io.ReadAll` 把文件读入 buffer，再用 `bytes.NewReader` 传给 HTTP body。

```go
// ❌ 错误：Doris 会返回 0 行
file, _ := os.Open(path)
req, _ := http.NewRequest("PUT", url, file)

// ✅ 正确：Content-Length 显式设置
data, _ := io.ReadAll(file)
req, _ := http.NewRequest("PUT", url, bytes.NewReader(data))
```

### 2. 必须显式提供 `Columns` 映射

调用 `ImportFile` 时**必须**传入 `LoadOptions.Columns`，其值与 Doris 表 `DESC` 输出列顺序一致。Doris 按位置匹配，避免列名差异或字段顺序不同导致的 0 行入库。

### 3. 写入器一步关闭

`parquet-go` 的 `WriteStop()` 内部调用 `pw.Close()`，自动完成：
- flush 内存中的 row group 到磁盘
- 写 parquet footer（schema、statistics）
- 关闭文件

**不要在 `WriteStop` 后再调 `Close()`**，会触发 "close of closed channel" panic。

## 故障排查

### 现象：HTTP 200 + Status=Success + 0 行入库

最常见的 3 个原因：

| 原因 | 检查方法 | 修复 |
|------|----------|------|
| HTTP 用了 chunked encoding | 用 `httputil.DumpRequestOut` 检查请求头是否有 `Transfer-Encoding: chunked` | 用 `io.ReadAll` 读入 buffer 再传 |
| `Columns` 与 Doris 表列顺序不一致 | 对比 `DESC` 输出与 `LoadOptions.Columns` 值 | 调整 Columns 顺序 |
| `Record` 字段类型与 Doris 列类型不匹配 | 比对两边的 schema | 按"Doris 表结构要求"小节调整 Go 类型 |

### 现象：HTTP 非 200 或 Status="Fail"

查看响应中的 `Message` 与 `ErrorURL`，`ErrorURL` 通常指向 Doris 的导入错误日志（如 `http://BE:8040/api/_load_error_log?file=...`）。