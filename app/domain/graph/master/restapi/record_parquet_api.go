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

// RecordParquetAPI
// @Description: 银行流水 parquet 导入 REST API
type RecordParquetAPI struct {
	env      *env.Env
	rootPath string
	service  *service.RecordParquetService
}

func NewRecordParquetAPI(env *env.Env, rootPath string) *RecordParquetAPI {
	return &RecordParquetAPI{
		env:      env,
		rootPath: rootPath,
		service:  service.NewRecordParquetService(env),
	}
}

// ImportRecordParquetPath 路径参数
type ImportRecordParquetPath struct {
	CaseId string `path:"caseId" required:"true"`
}

// ImportRecordParquetRequest 请求体
type ImportRecordParquetRequest struct {
	S3Path      string `json:"s3Path" validate:"required" desc:"s3://bucket/key 或 bucket/key"`
	BatchSize   int    `json:"batchSize,omitempty" desc:"批次大小, 默认 500"`
	Parallel    bool   `json:"parallel,omitempty" desc:"是否并行"`
	Concurrency int    `json:"concurrency,omitempty" desc:"并发数, 默认 1"`
	Retries     int    `json:"retries,omitempty" desc:"重试次数, 默认 1"`
}

// NewAPIController 注册路由
func (s *RecordParquetAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath, "master.RecordParquetAPI", s)
	ctl.Post("/case/{caseId}/record/parquet/import", "Import")
	ctl.Handle(iris.MethodOptions, "/case/{caseId}/record/parquet/import", "Check")
	return ctl
}

// Import 导入银行流水 parquet
func (s *RecordParquetAPI) Import(ctx context.Context, req *ImportRecordParquetRequest, path *ImportRecordParquetPath) (*dao.ImportResult, error) {
	opts := service.ImportOptions{
		BatchSize:   req.BatchSize,
		Parallel:    req.Parallel,
		Concurrency: req.Concurrency,
		Retries:     req.Retries,
	}
	return s.service.ImportFromS3(ctx, req.S3Path, path.CaseId, opts)
}

// Check 健康检查
func (s *RecordParquetAPI) Check(ctx context.Context) error {
	logs.Infofmt(ctx, "graph/master/restapi/record-parquet/import:check")
	return nil
}
