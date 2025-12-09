package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	client "github.com/ory/kratos-client-go"
)

type AuthAPI struct {
	env               *env.Env
	authService       *service.AuthService
	userService       *service.UserService
	tenantService     *service.TenantService
	tenantUserService *service.TenantUserService
	oryService        *service2.OryService
	jwtService        *service.JwtService
	rootPath          string
}

func NewAuthAPI(env *env.Env, rootPath string) *AuthAPI {
	return &AuthAPI{
		env:               env,
		authService:       service.NewAuthService(),
		userService:       service.NewUserService(),
		tenantService:     service.NewTenantService(),
		tenantUserService: service.NewTenantUserService(),
		oryService:        service2.NewOryService(),
		jwtService:        service.NewJwtService(),
		rootPath:          rootPath,
	}
}

func (s *AuthAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.authService = service.NewAuthService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.AuthAPI", s)
	ctl.Post("/login", "Login")
	ctl.Post("/generate", "Generate")
	return ctl
}

func (s *AuthAPI) Login(ctx context.Context, cmd *command.LoginCommand) (*model.LoginUser, error) {
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
	user, err := s.userService.FindUsingByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	loginUser := &model.LoginUser{
		Id:            user.Id,
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
		SessionId:     login.Session.Id,
	}

	return loginUser, nil
}

func (s *AuthAPI) Generate(ctx context.Context, cmd *command.GenerateJwtCommand) (*model.LoginResult, error) {
	user, err := s.userService.FindById(ctx, cmd.Data.UserId)
	if err != nil {
		return nil, err
	}

	loginUser := &model.LoginUser{
		Id:            user.Id,
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
		SessionId:     cmd.Data.SessionId,
	}

	tenant, err := s.tenantService.FindById(ctx, cmd.Data.TenantId)
	if err != nil {
		return nil, err
	}
	loginUser.TenantId = tenant.Id
	loginUser.TenantName = tenant.Name

	tu, err := s.tenantUserService.FindByTenantIdAndUserId(ctx, tenant.Id, user.Id)
	if err != nil {
		return nil, err
	}
	loginUser.IsAdmin = tu.IsAdmin

	session, err := s.oryService.GetSession(ctx, cmd.Data.SessionId)
	if err != nil {
		return nil, err
	}

	jwt, err := s.jwtService.Generate(&session.Session, loginUser)

	return &model.LoginResult{LoginSession: session, LoginJwt: jwt, LoginUser: *loginUser}, nil
}
