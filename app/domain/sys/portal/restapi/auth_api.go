package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	client "github.com/ory/kratos-client-go"
)

type AuthAPI struct {
	env         *env.Env
	authService *service.AuthService
	userService *service.UserService
	jwtService  *service.JwtService
	rootPath    string
}

func NewAuthAPI(env *env.Env, rootPath string) *AuthAPI {
	return &AuthAPI{
		env:         env,
		authService: service.NewAuthService(),
		userService: service.NewUserService(),
		jwtService:  service.NewJwtService(),
		rootPath:    rootPath,
	}
}

func (s *AuthAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.authService = service.NewAuthService()
	ctl := restapi.NewController(app, s.rootPath+"/auth", "sys.AuthAPI", s)
	ctl.Post("/login", "Login")
	return ctl
}

func (s *AuthAPI) Login(ctx context.Context, cmd *command.LoginCommand) (*model.LoginResult, error) {
	flow, err := s.authService.CreateLoginFlow(ctx)
	if err != nil {
		return nil, err
	}
	body := client.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &client.UpdateLoginFlowWithPasswordMethod{
			Password:   cmd.Data.Password,
			Identifier: cmd.Data.Identifier,
			Method:     cmd.Data.Method,
		},
	}
	login, err := s.authService.Login(ctx, flow.Id, body)
	if err != nil {
		if err.Error() == "400 Bad Request" {
			return nil, errors.New("错误的用户名或密码")
		}
		return nil, err
	}

	ident, ok := login.Session.Identity.Traits.(map[string]interface{})
	if !ok {
		return nil, errors.New("Traits转Map失败")
	}
	account, _ := ident["account"].(string)
	user, err := s.userService.FindByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	loginUser := &model.LoginUser{
		Id:            user.Id,
		TenantId:      "test",
		TenantName:    "test",
		Account:       user.Account,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		Address:       user.Address,
		Gender:        user.Gender,
		Work:          user.Work,
		HeadPicture:   user.HeadPicture,
		Status:        string(user.Status),
		OryIdentityId: user.OryIdentityId,
	}
	jwt, err := s.jwtService.Generate(&login.Session, loginUser)

	return &model.LoginResult{LoginSession: login, LoginJwt: jwt, LoginUser: *loginUser}, nil
}
