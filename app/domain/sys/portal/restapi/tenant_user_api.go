package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type TenantUserAPI struct {
	env               *env.Env
	tenantUserService *service.TenantUserService
	rootPath          string
}

func NewTenantUserAPI(env *env.Env, rootPath string) *TenantUserAPI {
	return &TenantUserAPI{
		env:               env,
		tenantUserService: service.NewTenantUserService(),
		rootPath:          rootPath,
	}
}

func (s *TenantUserAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tenantUserService = service.NewTenantUserService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.TenantUserAPI", s)
	ctl.Post("/tenant-user", "Create")
	ctl.Put("/tenant-user", "Update")
	ctl.Delete("/tenant-user:batch", "DeleteBatch", restapi.WithParamsInBody(true))
	return ctl
}

func (s *TenantUserAPI) Create(ctx context.Context, cmd *command.TenantUserCreateCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {

		err := s.tenantUserService.DeleteByTenantIdsAndUserIds(ctx, cmd.Data.TenIds, cmd.Data.UserIds)
		if err != nil {
			return err
		}

		tus := make([]*model.TenantUser, 0)
		for _, tenantId := range cmd.Data.TenIds {
			for _, userId := range cmd.Data.UserIds {
				tu, _ := model.NewTenantUser()
				tu.Id = idutils.NewId()
				tu.TenId = tenantId
				tu.UserId = userId
				tus = append(tus, tu)
			}
		}

		return s.tenantUserService.CreateMany(ctx, tus)
	})
}

func (s *TenantUserAPI) Update(ctx context.Context, cmd *command.TenantUserUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"is_admin"})
	return s.tenantUserService.Update(ctx, &cmd.Data, opts)
}

func (s *TenantUserAPI) DeleteBatch(ctx context.Context, cmd *command.TenantUserDeleteBatchCommand) error {
	return s.tenantUserService.DeleteBatch(ctx, cmd.Data.Ids)
}
