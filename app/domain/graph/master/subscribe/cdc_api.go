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
	recordService *service.RagRecordService
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

func (s *CdcAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath, "CdcAPI", s)
	ctl.Handle(iris.MethodPost, "/master-cdc-graph", "DataChange")
	ctl.Handle(iris.MethodOptions, "/master-cdc-graph", "DaprOptions")
	return ctl
}

func (s *CdcAPI) DaprOptions(ctx context.Context) error {
	return nil
}

func (s *CdcAPI) DataChange(ctx context.Context, cdc *model.CDCRecord) error {
	if cdc.Table == "record" {
		return s.recordService.DataChange(ctx, cdc)
	}
	return s.masterService.DataChange(ctx, cdc)
}
