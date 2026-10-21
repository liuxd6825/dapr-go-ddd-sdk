package subscribe

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	graphconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/miniofs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/spf13/afero"
)

// RecordEventSubHandler
// @Description: 处理数据导入事件，触发 s3 parquet 批量导入 Neo4j
type RecordEventSubHandler struct {
	rootPath      string
	env           *envs.Env
	fs            afero.Fs
	recordService *service.RecordService
}

// NewRecordEventHandler
//
//	@Description:
//	@param env
//	@param baseUrl
//	@param s3Key  "default",
//	@param s3bucketKey "importData"
//	@return *RecordEventSubHandler
//	@return error
func NewRecordEventHandler(env *envs.Env, baseUrl string, s3Key string, s3bucketKey string, graphCfg graphconfigs.Config) (*RecordEventSubHandler, error) {
	fs, err := miniofs.NewFs(env, s3Key, s3bucketKey)
	if err != nil {
		return nil, err
	}
	return &RecordEventSubHandler{
		rootPath:      baseUrl,
		env:           env,
		fs:            fs,
		recordService: service.NewRecordServiceWithEnv(env, graphCfg),
	}, nil
}

func (s *RecordEventSubHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, "subscribe/master/event", "RecordEventHandler", s)
	ctl.EventHandle("record-import-master-event", "RecordImportMasterEvent")
	ctl.Handle(iris.MethodOptions, "record-import-master-event", "Check")
	return ctl
}

func (s *RecordEventSubHandler) Check(ctx context.Context) error {
	logs.Infofmt(ctx, "graph/master/subscribe/import/record-import-master-event:check")
	return nil
}

// RecordImportMasterEvent 领域事件入口
// 接收 RecordImportMasterEvent，调用 RecordParquetService.ImportFromS3 触发按日期聚合的批量导入。
func (s *RecordEventSubHandler) RecordImportMasterEvent(ctx context.Context, ev *event.RecordImportMasterEvent) error {
	if ev == nil || ev.Data == nil {
		return errors.New("event or event.Data is nil")
	}
	data := ev.Data

	s3DataUrl := data.S3FileUrl
	caseId := data.CaseId
	taskId := data.TaskId

	if s3DataUrl == "" {
		return errors.New("event.Data.S3FileUrl is empty")
	}
	if caseId == "" {
		return errors.New("event.Data.CaseId is empty")
	}
	if taskId == "" {
		return errors.New("event.Data.TaskId is empty")
	}

	// tenantId 优先从 ctx 拿 (订阅路径已注入), 事件自身作为兜底
	tenantId := appctx.GetTenantId2(ctx)
	if tenantId == "" {
		tenantId = ev.TenantId
	}
	if tenantId == "" {
		return errors.New("tenantId is empty in context and event")
	}

	logs.InfoMsg(ctx, "RecordImportMasterEvent start",
		" eventId=", ev.Id,
		" occurredOn=", ev.CreatedTime.Format(time.RFC3339Nano),
		" tenantId=", tenantId,
		" caseId=", caseId,
		" s3URL=", s3DataUrl,
		" taskId=", data.TaskId,
		" fileId=", data.FileId,
		" fileName=", data.FileName,
	)

	s3SparkJarURL := ""

	sparkTaskId, err := s.recordService.CreateImportTask(ctx, caseId, taskId, s3DataUrl, s3SparkJarURL, taskId)
	if err != nil {
		logs.ErrorErr(ctx, err)
		logs.Errorfmt(ctx, "RecordImportMasterEvent failed: eventId=%s caseId=%s err=%v", ev.Id, caseId, err)
		return err
	}

	logs.Infofmt(ctx, "RecordImportMasterEvent end Spark Task ID: %s", sparkTaskId)

	return nil
}
