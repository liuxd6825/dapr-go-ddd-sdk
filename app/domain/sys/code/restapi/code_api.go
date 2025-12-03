package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CodeAPI struct {
	env         *env.Env
	codeService *service.CodeService
	rootPath    string
}

func NewCodeAPI(env *env.Env, rootPath string) *CodeAPI {
	return &CodeAPI{
		env:         env,
		codeService: service.NewCodeService(),
		rootPath:    rootPath,
	}
}

func (s *CodeAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.CodeAPI", s)
	ctl.GetOne("/code:new", "NewCode")
	return ctl
}

func (s *CodeAPI) NewCode(ctx context.Context, newCmd *command.CodeNewCommand) (*command.CodeNewCommandResult, error) {
	code, err := s.codeService.New(ctx, newCmd)
	if err != nil {
		return nil, err
	}
	return &command.CodeNewCommandResult{Code: code}, nil
}
