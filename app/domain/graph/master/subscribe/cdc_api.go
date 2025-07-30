package subscribe

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CdcAPI struct {
	env           *env.Env
	rootPath      string
	masterService *service.MasterService
	queryService  *service.QueryService
	recordService *service.RecordService
}

func NewCdcAPI(env *env.Env, rootPath string) *CdcAPI {
	masterService := service.NewMasterService().Init()
	return &CdcAPI{
		env:           env,
		rootPath:      rootPath,
		masterService: masterService,
		recordService: service.NewRecordService(),
	}
}

func (s *CdcAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath, s)
	ctl.Handle(iris.MethodPost, "/master-cdc-graph", "DataChange")
	ctl.Handle(iris.MethodOptions, "/master-cdc-graph", "DaprOptions")
	return nil
}

func (s *CdcAPI) DaprOptions(ctx context.Context) error {
	return nil
}

func (s *CdcAPI) DataChange(ctx context.Context, record *model.CDCRecord) error {
	if record.Table == "record" {
		return s.recordService.DataChange(ctx, record)
	}
	return s.masterService.DataChange(ctx, record)
}
