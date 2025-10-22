package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	client "github.com/ory/kratos-client-go"
)

type UserAPI struct {
	env         *env.Env
	userService *service.UserService
	oryService  *service2.OryService
	rootPath    string
}

func NewUserAPI(env *env.Env, rootPath string) *UserAPI {
	return &UserAPI{
		env:         env,
		userService: service.NewUserService(),
		oryService:  service2.NewOryService(),
		rootPath:    rootPath,
	}
}

func (s *UserAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.userService = service.NewUserService()
	ctl := restapi.NewController(app, s.rootPath+"/portal", "sys.UserAPI", s)
	ctl.Post("/user", "Create")
	ctl.Put("/user", "Update")
	ctl.Put("/user:reset-password", "ResetPassword")
	ctl.Put("/user:update-password", "UpdatePassword")
	ctl.Delete("/user", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/user/{id}", "FindById")
	ctl.GetPaging("/user", "FindPaging")
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
	cmd.UpdateMask = []string{"name", "phone", "email", "address", "gender", "work"}
	return s.userService.Update(ctx, cmd)
}

func (s *UserAPI) UpdatePassword(ctx context.Context, cmd *command.UpdatePasswordCommand) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		if cmd.Data.Password != cmd.Data.ConfirmPassword {
			return errors.New("密码不一致")
		}
		user, err := s.userService.FindByAccount(ctx, cmd.Data.Account)
		if err != nil {
			return err
		}
		password := cmd.Data.Password
		updateIdentityBody := client.UpdateIdentityBody{
			Credentials: &client.IdentityWithCredentials{
				Password: &client.IdentityWithCredentialsPassword{
					Config: &client.IdentityWithCredentialsPasswordConfig{
						Password: &password,
					},
				},
			},
			Traits: map[string]interface{}{
				"account": user.Account,
				"email":   user.Email,
			},
		}
		_, err = s.oryService.UpdateIdentity(ctx, user.OryIdentityId, updateIdentityBody)
		if err != nil {
			return err
		}
		data, _ := model.NewUser()
		data.Id = user.Id
		data.Password = password
		return s.userService.UpdateData(ctx, data, idao.NewCallOptions().SetUpdateFields([]string{"password"}))
	})
	return err
}

func (s *UserAPI) ResetPassword(ctx context.Context, cmd *command.UserUpdateCommand) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		password := "123@@abc"
		updateIdentityBody := client.UpdateIdentityBody{
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
		_, err := s.oryService.UpdateIdentity(ctx, cmd.Data.OryIdentityId, updateIdentityBody)
		if err != nil {
			return err
		}
		cmd.UpdateMask = []string{"password"}
		cmd.Data.Password = password
		return s.userService.Update(ctx, cmd)
	})
	return err
}

func (s *UserAPI) Delete(ctx context.Context, cmd *command.UserDeleteCommand) error {
	err := tx.StartTx(ctx, []string{s.userService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		err := s.oryService.DeleteIdentity(ctx, cmd.Data.OryIdentityId)
		if err != nil {
			return err
		}
		return s.userService.Delete(ctx, cmd)
	})
	return err
}

func (s *UserAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.User, error) {
	return s.userService.FindById(ctx, qry)
}

func (s *UserAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.User], error) {
	return s.userService.FindPaging(ctx, qry)
}
