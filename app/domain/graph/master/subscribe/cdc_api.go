package subscribe

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CdcAPI struct {
	env           *env.Env
	rootPath      string
	masterService *service.MasterService
	queryService  *service.QueryService
	recordService *service.RagRecordService
}

var dataChangeURL = "/subscribe/graph/master/cdc/data-change"

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
	ctl := restapi.NewController(app, "", "CdcAPI", s)
	ctl.CDCHandle(dataChangeURL, "DataChange")
	ctl.Handle(iris.MethodOptions, dataChangeURL, "DaprOptions")
	return ctl
}

func (s *CdcAPI) DaprOptions(ctx context.Context) error {
	logs.Infofmt(ctx, "dapr subscribe %s ", dataChangeURL)
	return nil
}

func (s *CdcAPI) DataChange(ctx context.Context, ictx iris.Context, cdc *restapi.CDCRecord) error {
	logs.Infofmt(ctx, "%s CdcAPI.DataChange(CDCRecord) table:[%s]", dataChangeURL, cdc.Table)
	if cdc.Table == "record" {
		return s.recordService.DataChange(ctx, cdc)
	}
	return s.masterService.DataChange(ctx, cdc)
}
