package subscribe

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CdcAPI struct {
	env           *envs.Env
	rootPath      string
	masterService *service.MasterService
	queryService  *service.MasterQueryService
}

var dataChangeURL = "/subscribe/graph/master/cdc/data-change"

func NewCdcAPI(env *envs.Env, rootPath string, hugeDbKey string) *CdcAPI {
	hugeClient := env.GetHugeClient(hugeDbKey)
	masterService := service.NewMasterService(hugeDbKey, hugeClient).Init()
	return &CdcAPI{
		env:           env,
		rootPath:      rootPath,
		masterService: masterService,
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

func (s *CdcAPI) DataChange(ctx context.Context, ictx iris.Context, cdc *dbevent.CDCRecord) error {
	logs.Infofmt(ctx, "%s CdcAPI.DataChange(CDCRecord) table:[%s]", dataChangeURL, cdc.Table)
	return s.masterService.DataChange(ctx, cdc)
}
