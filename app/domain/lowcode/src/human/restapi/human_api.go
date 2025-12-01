package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/lowcode/src/human/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type HumanAPI struct {
	env          *env.Env
	rootPath     string
	humanService *service.HumanService
}

func NewHumanAPI(env *env.Env, rootPath string) *HumanAPI {
	return &HumanAPI{
		env:          env,
		rootPath:     rootPath,
		humanService: service.NewHumanService(),
	}
}

func (s *HumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath, "HumanApi", s)
	ctl.Post("human", "Create")
	ctl.Put("human", "Update")
	ctl.Delete("human", "Delete")
	ctl.GetPaging("human", "FindPaging")
	return ctl
}

func (s *HumanAPI) Create(ctx context.Context, cmd *command.HumanCreateCommand) error {
	return s.humanService.Create(ctx, cmd)
}

func (s *HumanAPI) Update(ctx context.Context, cmd *command.HumanCreateCommand) error {
	return s.humanService.Create(ctx, cmd)
}
