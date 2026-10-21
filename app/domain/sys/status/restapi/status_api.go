package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type StatusApi struct {
	service  *service.StatusService
	env      *envs.Env
	rootPath string
}

func NewStatusApi(env *envs.Env, rootPath string, mongoDbKey string) *StatusApi {
	statusService := service.NewStatusService(mongoDbKey)
	return &StatusApi{
		env:      env,
		rootPath: rootPath,
		service:  statusService,
	}
}

func (s *StatusApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/sys/status", "sys.status", s.service)
	controller.Post("", "Create")
	controller.Put("", "Update")
	controller.Delete("", "Delete")
	controller.GetOne("/{id}", "FindById")
	controller.GetPaging("", "FindPaging")
	controller.GetData(":user-id", "FindByUserId")
	controller.GetData(":all", "FindAll")
	return controller
}
