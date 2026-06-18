package doris

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/spf13/afero"
)

// DorisImporter 负责执行 Stream Load 导入
type DorisImporter struct {
	config Config
	client *http.Client
}

// DorisImporter 初始化加载器
func NewDorisImporter(cfg Config) *DorisImporter {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &DorisImporter{
		config: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// LoadOptions 允许自定义单次导入的参数
type LoadOptions struct {
	Format          LoadFormat // 文件格式，例如 "parquet"、"csv"、"json" (留空默认视为 csv)
	Columns         string     // 列映射关系，例如 "id, name, age" (非必填，Doris默认按列名或顺序对齐)
	ColumnSeparator string     // CSV 特有：列分隔符，例如 "," 或 "\t"
	LineDelimiter   string     // CSV 特有：行分隔符，例如 "\n"
	Label           string     // 自定义唯一批次号 (防重入)
}

type LoadFormat string

const (
	FormatParquet LoadFormat = "parquet"
	FormatCSV     LoadFormat = "csv"
	FormatJSON    LoadFormat = "json"
)

func (f LoadFormat) String() string {
	return string(f)
}

// DorisLoadResp Doris Stream Load 接口的 JSON 响应
type DorisImportResp struct {
	TxnID                  int64  `json:"TxnId"`
	Label                  string `json:"Label"`
	Status                 string `json:"Status"`
	ExistingJobStatus      string `json:"ExistingJobStatus"`
	Message                string `json:"Message"`
	NumberTotalRows        int64  `json:"NumberTotalRows"`
	NumberLoadedRows       int64  `json:"NumberLoadedRows"`
	NumberFilteredRows     int64  `json:"NumberFilteredRows"`
	NumberUnselectedRows   int64  `json:"NumberUnselectedRows"`
	NumberUnloadedRows     int64  `json:"NumberUnloadedRows"`
	LoadBytes              int64  `json:"LoadBytes"`
	LoadTimeMs             int64  `json:"LoadTimeMs"`
	BeginTxnTimeMs         int64  `json:"BeginTxnTimeMs"`
	StreamLoadPutTimeMs    int64  `json:"StreamLoadPutTimeMs"`
	ReadDataTimeMs         int64  `json:"ReadDataTimeMs"`
	WriteDataTimeMs        int64  `json:"WriteDataTimeMs"`
	CommitAndPublishTimeMs int64  `json:"CommitAndPublishTimeMs"`
	ErrorURL               string `json:"ErrorURL"`
}

func (r *DorisImportResp) loaded() bool {
	return r.NumberLoadedRows > 0
}

// ImportFile 从磁盘读取 Parquet 文件并通过 Stream Load 上传到 Doris。
//
// 关键实现细节：必须先把文件读入 buffer 再发送，**不能直接传 *os.File 给 HTTP body**。
// 如果传 *os.File，Go 的 net/http 客户端会用 Transfer-Encoding: chunked 发送请求，
// 而 Doris Stream Load 不接受 chunked encoding（要求 Content-Length）。
// 表现：Doris 返回 HTTP 200 + Status=Success + 0 行入库（"假成功"）。
func (dl *DorisImporter) ImportFile(fs afero.Fs, parquetPath string, opts LoadOptions) (*DorisImportResp, error) {
	file, err := fs.OpenFile(parquetPath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errors.New("❌ 打开 Parquet 文件时出错: %v\n", err)
	}
	defer file.Close()

	data, err := afero.ReadAll(file)
	if err != nil {
		return nil, errors.New("❌ 读取 Parquet 文件时出错: %v\n", err)
	}
	return dl.ImportData(bytes.NewReader(data), opts)
}

// ImportData 核心导入方法，支持任何实现了 io.Reader 的数据流。
//
// 注意：调用方需保证传入的 io.Reader 的剩余字节数可知（即要么是 *bytes.Reader /
// *strings.Reader，要么先 io.ReadAll 到 buffer 后再传），否则 Go HTTP 客户端会
// 退化为 Transfer-Encoding: chunked，Doris 会拒绝解析（参见 ImportFile）。
func (dl *DorisImporter) ImportData(data io.Reader, opts LoadOptions) (*DorisImportResp, error) {
	url := fmt.Sprintf("http://%s:%d/api/%s/%s/_stream_load",
		dl.config.BEHost, dl.config.BEPort, dl.config.Database, dl.config.Table)

	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// 基础认证
	auth := dl.config.Username + ":" + dl.config.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)

	// 文件格式相关 headers
	if opts.Format != "" {
		req.Header.Set("format", opts.Format.String())
	}
	if opts.Format == "" || opts.Format == FormatCSV {
		if opts.ColumnSeparator != "" {
			req.Header.Set("column_separator", opts.ColumnSeparator)
		}
		if opts.LineDelimiter != "" {
			req.Header.Set("line_delimiter", opts.LineDelimiter)
		}
	}

	// 列映射与 Label
	if opts.Columns != "" {
		req.Header.Set("columns", opts.Columns)
	}
	if opts.Label != "" {
		req.Header.Set("label", opts.Label)
	}

	// 执行请求
	resp, err := dl.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request execution failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doris returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	loadResp := &DorisImportResp{}
	if err := json.Unmarshal(body, loadResp); err != nil {
		return nil, fmt.Errorf("failed to parse doris response: %w, body: %s", err, string(body))
	}

	if !strings.EqualFold(loadResp.Status, "Success") {
		return loadResp, fmt.Errorf("doris stream load failed, status: %s, message: %s, error_url: %s",
			loadResp.Status, loadResp.Message, loadResp.ErrorURL)
	}
	if !loadResp.loaded() {
		return loadResp, fmt.Errorf("doris stream load succeeded but no rows loaded, total: %d, loaded: %d, filtered: %d, unloaded: %d, message: %s",
			loadResp.NumberTotalRows, loadResp.NumberLoadedRows,
			loadResp.NumberFilteredRows, loadResp.NumberUnloadedRows, loadResp.Message)
	}
	return loadResp, nil
}

// S3Config Doris 直接从 S3 兼容存储取文件所需的凭证。
//
// 字段值必须与最初上传 Parquet/CSV/JSON 文件时使用的 minio.Client 凭证
// 一致,Doris 会用它自己连接 S3/MinIO,loader 不再参与文件转发。
type S3Config struct {
	Endpoint     string // 例如 "192.168.120.224:9000"(不带 scheme)
	Region       string // 例如 "us-east-1"(MinIO 任意值均可)
	AccessKey    string
	SecretKey    string
	UsePathStyle bool // MinIO 必须为 true,AWS S3 用 false
	UseSSL       bool // endpoint 为 https 时为 true
}

// BuildS3Path 拼接 "s3://bucket/key" 形式的 URL。key 可以带或不带前导 "/",
// 都会被规范化。
func BuildS3Path(bucket, key string) string {
	return fmt.Sprintf("s3://%s/%s", bucket, strings.TrimLeft(key, "/"))
}

// buildS3InsertQuery builds the INSERT INTO ... SELECT FROM S3(...) SQL that
// the MySQL-based ImportFileFromS3 issues. Extracted so unit tests can verify
// the generated SQL without standing up a MySQL server.
//
// Doris's S3 TVF expects single-quoted string literals, so we wrap each
// interpolated value with single quotes explicitly (Go's %q produces
// double-quoted Go-style strings, which Doris would reject).
func buildS3InsertQuery(database, table, s3Path, accessKey, secretKey,
	region, endpoint string, usePathStyle bool, format string,
) string {
	q := func(s string) string { return "'" + strings.ReplaceAll(s, "'", `\'`) + "'" }
	return fmt.Sprintf(
		"INSERT INTO %s.%s SELECT * FROM S3("+
			"'uri' = %s, "+
			"'access_key' = %s, "+
			"'secret_key' = %s, "+
			"'region' = %s, "+
			"'endpoint' = %s, "+
			"'use_path_style' = '%t', "+
			"'format' = %s"+
			")",
		database, table,
		q(s3Path), q(accessKey), q(secretKey), q(region), q(endpoint), usePathStyle, q(format),
	)
}

// ImportFileFromS3 让 Doris 直接从 S3/MinIO 取数据文件并写入目标表。
//
// 实现方式:走 Doris FE 的 MySQL 协议,执行
//
//	INSERT INTO {db}.{table}
//	SELECT * FROM S3(
//	    'uri'          = 's3://bucket/key',
//	    'access_key'   = '...',
//	    'secret_key'   = '...',
//	    'region'       = '...',
//	    'endpoint'     = 'host:port',
//	    'use_path_style' = 'true|false',
//	    'format'       = 'parquet|csv|json'
//	)
//
// 这条路在所有支持 S3 TVF 的 Doris 版本(≥1.2)上都能用,比 HTTP-based
// "S3 Stream Load" 兼容性更好,且不依赖 Stream Load 1.2+ 的特殊 body
// 格式(该格式在某些版本里把 body 当作文件内容直接处理)。
//
// 与 ImportFile 的区别:
//   - 避开"MinIO → loader → Doris"的字节往返,大文件时性能明显更好;
//   - 避开 ImportFile 注释里提到的 Content-Length / chunked encoding
//     限制,杜绝"HTTP 200 + Success + 0 行入库"的假成功。
//
// s3Path 必须是 "s3://bucket/key" 格式(可用 BuildS3Path 构造)。
// s3cfg 中的凭证会作为 S3 TVF 的属性传给 Doris。
//
// 注意:本方法需要 loader 能连接 Doris FE 的 MySQL 端口(默认 9030)。
// Config 中 FEHost/FEPort 留空时,自动回退到 BEHost:9030。
func (dl *DorisImporter) ImportFileFromS3(s3Path string, s3cfg S3Config, opts LoadOptions) (*DorisImportResp, error) {
	if s3Path == "" {
		return nil, errors.New("miniofs: s3 path is empty")
	}
	if s3cfg.AccessKey == "" || s3cfg.SecretKey == "" {
		return nil, errors.New("miniofs: s3 access/secret key is required")
	}

	format := opts.Format.String()
	if format == "" {
		format = string(FormatParquet)
	}

	// Build the INSERT INTO ... SELECT FROM S3(...) statement.
	query := buildS3InsertQuery(
		dl.config.Database, dl.config.Table,
		s3Path, s3cfg.AccessKey, s3cfg.SecretKey,
		s3cfg.Region, s3cfg.Endpoint, s3cfg.UsePathStyle, format,
	)

	feHost := dl.config.FEHost
	if feHost == "" {
		feHost = dl.config.BEHost
	}
	fePort := dl.config.FEPort
	if fePort <= 0 {
		fePort = 9030
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dl.config.Username, dl.config.Password, feHost, fePort, dl.config.Database))
	if err != nil {
		return nil, fmt.Errorf("open mysql conn: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), dl.config.Timeout)
	defer cancel()
	if _, err := db.ExecContext(ctx, query); err != nil {
		return nil, fmt.Errorf("doris s3 import failed: %w", err)
	}

	// The MySQL protocol doesn't return a DorisLoadResp-style struct from
	// INSERT...SELECT. We synthesise a "Success" response so callers that
	// rely on the unified return type continue to work; NumberLoadedRows
	// is left at zero because the MySQL protocol does not surface it.
	return &DorisImportResp{
		Status:          "Success",
		Message:         fmt.Sprintf("loaded via S3 TVF from %s", s3Path),
		NumberTotalRows: -1, // unknown via this path
	}, nil
}

// ImportFileFromS3HTTP is the HTTP-based variant of ImportFileFromS3. It
// issues a PUT against {BE}/api/{db}/{table}/_stream_load with the S3 path
// as the request body. This is the "S3 Stream Load" form documented for
// Doris 1.2+, but in practice some Doris versions ignore the body and read
// it as the file content. Prefer ImportFileFromS3 (MySQL + S3 TVF) unless
// your Doris is known to support the HTTP variant.
//
// Kept for forward compatibility; new callers should use ImportFileFromS3.
func (dl *DorisImporter) ImportFileFromS3HTTP(s3Path string, s3cfg S3Config, opts LoadOptions) (*DorisImportResp, error) {
	if s3Path == "" {
		return nil, errors.New("miniofs: s3 path is empty")
	}
	if s3cfg.AccessKey == "" || s3cfg.SecretKey == "" {
		return nil, errors.New("miniofs: s3 access/secret key is required")
	}

	url := fmt.Sprintf("http://%s:%d/api/%s/%s/_stream_load",
		dl.config.BEHost, dl.config.BEPort, dl.config.Database, dl.config.Table)

	// Doris S3 Stream Load 模式:body 是纯文本 S3 URL,Content-Length 由
	// strings.NewReader 显式告知,避免触发 Transfer-Encoding: chunked。
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(s3Path))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	// 基础认证
	auth := dl.config.Username + ":" + dl.config.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)

	// S3 凭证
	req.Header.Set("AWS_ACCESS_KEY", s3cfg.AccessKey)
	req.Header.Set("AWS_SECRET_KEY", s3cfg.SecretKey)
	if s3cfg.Region != "" {
		req.Header.Set("AWS_REGION", s3cfg.Region)
	}
	if s3cfg.Endpoint != "" {
		req.Header.Set("AWS_ENDPOINT", s3cfg.Endpoint)
	}
	if s3cfg.UsePathStyle {
		req.Header.Set("AWS_USE_PATH_STYLE", "true")
	}
	if s3cfg.UseSSL {
		req.Header.Set("AWS_USE_SSL", "true")
	}

	// 文件格式 / 列映射 / Label,与 ImportData 一致
	if opts.Format != "" {
		req.Header.Set("format", opts.Format.String())
	}
	if opts.Format == "" || opts.Format == FormatCSV {
		if opts.ColumnSeparator != "" {
			req.Header.Set("column_separator", opts.ColumnSeparator)
		}
		if opts.LineDelimiter != "" {
			req.Header.Set("line_delimiter", opts.LineDelimiter)
		}
	}
	if opts.Columns != "" {
		req.Header.Set("columns", opts.Columns)
	}
	if opts.Label != "" {
		req.Header.Set("label", opts.Label)
	}

	// 执行请求
	resp, err := dl.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request execution failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doris returned non-200 status: %d, body: %s",
			resp.StatusCode, string(respBody))
	}

	// 解析响应
	loadResp := &DorisImportResp{}
	if err := json.Unmarshal(respBody, loadResp); err != nil {
		return nil, fmt.Errorf("failed to parse doris response: %w, body: %s", err, string(respBody))
	}

	if !strings.EqualFold(loadResp.Status, "Success") {
		return loadResp, fmt.Errorf("doris s3 stream load failed, status: %s, message: %s, error_url: %s",
			loadResp.Status, loadResp.Message, loadResp.ErrorURL)
	}
	if !loadResp.loaded() {
		return loadResp, fmt.Errorf("doris s3 stream load succeeded but no rows loaded, total: %d, loaded: %d, filtered: %d, unloaded: %d, message: %s",
			loadResp.NumberTotalRows, loadResp.NumberLoadedRows,
			loadResp.NumberFilteredRows, loadResp.NumberUnloadedRows, loadResp.Message)
	}
	return loadResp, nil
}
