package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	client "github.com/ory/kratos-client-go"
)

type AuthAPI struct {
	env         *env.Env
	authService *service.AuthService
	rootPath    string
}

func NewAuthAPI(env *env.Env, rootPath string) *AuthAPI {
	return &AuthAPI{
		env:         env,
		authService: service.NewAuthService(),
		rootPath:    rootPath,
	}
}

func (s *AuthAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.authService = service.NewAuthService()
	ctl := restapi.NewController(app, s.rootPath+"/auth", "sys.AuthAPI", s)
	ctl.GetData("/login-flow", "CreateLoginFlow")
	ctl.Post("/login", "Login")
	return ctl
}

func (s *AuthAPI) CreateLoginFlow(ctx context.Context) (*client.LoginFlow, error) {
	return s.authService.CreateLoginFlow(ctx)
}

func (s *AuthAPI) Login(ctx context.Context, cmd *command.LoginCommand) (*client.SuccessfulNativeLogin, error) {
	body := client.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &client.UpdateLoginFlowWithPasswordMethod{
			Password:   cmd.Data.Password,
			Identifier: cmd.Data.Identifier,
			Method:     cmd.Data.Method,
		},
	}
	login, err := s.authService.Login(ctx, cmd.Data.Flow, body)
	if err != nil {
		if err.Error() == "400 Bad Request" {
			return nil, errors.New("错误的用户名或密码")
		}
		return nil, err
	}
	return login, nil
}
