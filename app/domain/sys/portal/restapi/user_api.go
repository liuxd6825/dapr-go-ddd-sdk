package restapi

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	client "github.com/ory/kratos-client-go"
)

type UserAPI struct {
	env               *env.Env
	userService       *service.UserService
	oryService        *service2.OryService
	tenantUserService *service.TenantUserService
	rootPath          string
}

func NewUserAPI(env *env.Env, rootPath string) *UserAPI {
	return &UserAPI{
		env:               env,
		userService:       service.NewUserService(),
		oryService:        service2.NewOryService(),
		tenantUserService: service.NewTenantUserService(),
		rootPath:          rootPath,
	}
}

func (s *UserAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.userService = service.NewUserService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.UserAPI", s)
	ctl.Post("/user", "Create")
	ctl.Put("/user", "Update")
	ctl.Put("/user:reset-password", "ResetPassword")
	ctl.Put("/user:update-password", "UpdatePassword")
	ctl.Delete("/user", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/user:ident", "DeleteIdentity", restapi.WithParamsInBody(true))
	ctl.GetOne("/user/{id}", "FindById")
	ctl.GetPaging("/user", "FindPaging")
	ctl.GetPaging("/user:view", "FindPagingByTenantId")
	ctl.GetPaging("/user:out", "FindPagingByNotInTenantUser")
	return ctl
}

func (s *UserAPI) Create(ctx context.Context, cmd *command.UserCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		password := "123@@abc"
		body := client.CreateIdentityBody{
			Credentials: &client.IdentityWithCredentials{
				Password: &client.IdentityWithCredentialsPassword{
					Config: &client.IdentityWithCredentialsPasswordConfig{
						Password: &password,
					},
				},
			},
			Traits: map[string]interface{}{
				"account": cmd.Data.Account,
				"email":   cmd.Data.Email,
			},
		}
		ident, err := s.oryService.CreateIdentityExecute(ctx, body)
		if err != nil {
			if err.Error() == "409 Conflict" {
				return errors.New("账号或邮箱已存在")
			}
			return err
		}
		cmd.Data.OryIdentityId = ident.Id
		cmd.Data.Password = password
		return s.userService.Create(ctx, cmd)
	})
	return err
}

func (s *UserAPI) Update(ctx context.Context, cmd *command.UserUpdateCommand) error {
	return s.userService.Update(ctx, cmd)
}

func (s *UserAPI) UpdatePassword(ctx context.Context, cmd *command.UpdatePasswordCommand) error {
	if cmd.Data.Password != cmd.Data.ConfirmPassword {
		return errors.New("密码不一致")
	}
	user, err := s.userService.FindUsingByAccount(ctx, cmd.Data.Account)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("账号已禁用")
	}
	user.Password = cmd.Data.Password
	return s.updatePassword(ctx, user)
}

func (s *UserAPI) ResetPassword(ctx context.Context, cmd *command.UserUpdateCommand) error {
	cmd.Data.Password = "123@@abc"
	return s.updatePassword(ctx, &cmd.Data)
}

func (s *UserAPI) updatePassword(ctx context.Context, user *model.User) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		updateIdentityBody := client.UpdateIdentityBody{
			Credentials: &client.IdentityWithCredentials{
				Password: &client.IdentityWithCredentialsPassword{
					Config: &client.IdentityWithCredentialsPasswordConfig{
						Password: &user.Password,
					},
				},
			},
			Traits: map[string]interface{}{
				"account": user.Account,
				"email":   user.Email,
			},
		}
		_, err := s.oryService.UpdateIdentity(ctx, user.OryIdentityId, updateIdentityBody)
		if err != nil {
			return err
		}

		return s.userService.UpdateData(ctx, user, idao.NewCallOptions().SetUpdateFields([]string{"password"}))
	})
	return err
}

func (s *UserAPI) Delete(ctx context.Context, cmd *command.UserDeleteCommand) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		err := s.oryService.DeleteIdentity(ctx, cmd.Data.OryIdentityId)
		if err != nil {
			return err
		}
		err = s.tenantUserService.DeleteByUserId(ctx, cmd.Data.Id)
		if err != nil {
			return err
		}
		return s.userService.Delete(ctx, cmd)
	})
	return err
}

func (s *UserAPI) DeleteIdentity(ctx context.Context, cmd *command.UserDeleteIdentityCommand) error {
	return s.oryService.DeleteIdentity(ctx, cmd.Data.Id)
}

func (s *UserAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.User, error) {
	return s.userService.FindById(ctx, qry.Id)
}

func (s *UserAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.User], error) {
	ctx, _ = restapp.NewTestContext(context.Background(), func(option *restapp.ContextOption) {
		tenantId := service.SystemTenantId
		option.TenantId = &tenantId
	})
	return s.userService.FindPaging(ctx, qry)
}

func (s *UserAPI) FindPagingByTenantId(ctx context.Context, qry *query.FindPagingByTenantIdRequest) (idao.FindPagingResult[*model.UserView], error) {
	ctx, _ = restapp.NewTestContext(context.Background(), func(option *restapp.ContextOption) {
		tenantId := service.SystemTenantId
		option.TenantId = &tenantId
	})
	return s.userService.FindPagingByTenantId(ctx, qry.TenantId, qry)
}

func (s *UserAPI) FindPagingByNotInTenantUser(ctx context.Context, qry *query.FindPagingByTenantIdRequest) (idao.FindPagingResult[*model.User], error) {
	ctx, _ = restapp.NewTestContext(context.Background(), func(option *restapp.ContextOption) {
		tenantId := service.SystemTenantId
		option.TenantId = &tenantId
	})
	return s.userService.FindPagingByNotInTenantUser(ctx, qry.TenantId, qry)
}
