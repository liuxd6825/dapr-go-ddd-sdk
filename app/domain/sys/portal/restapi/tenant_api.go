package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	service4 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/service"
	service3 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	service5 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type TenantAPI struct {
	env               *env.Env
	tenantService     *service.TenantService
	userService       *service.UserService
	tenantUserService *service.TenantUserService
	oryService        *service2.OryService
	currencyService   *service3.CurrencyService
	bankService       *service4.BankService
	tagService        *service5.TagService
	rootPath          string
}

func NewTenantAPI(env *env.Env, rootPath string) *TenantAPI {
	return &TenantAPI{
		env:               env,
		tenantService:     service.NewTenantService(),
		userService:       service.NewUserService(),
		tenantUserService: service.NewTenantUserService(),
		oryService:        service2.NewOryService(),
		currencyService:   service3.NewCurrencyService(),
		bankService:       service4.NewBankService(),
		tagService:        service5.NewTagService(),
		rootPath:          rootPath,
	}
}

func (s *TenantAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tenantService = service.NewTenantService()
	ctl := restapi.NewController(app, s.rootPath+"/sys", "sys.TenantAPI", s)
	ctl.Post("/tenant", "Create")
	ctl.Post("/tenant:init", "InitTenant")
	ctl.Put("/tenant", "Update")
	ctl.Delete("/tenant", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/tenant/{id}", "FindById")
	ctl.GetPaging("/tenant", "FindPaging")
	ctl.GetData("/tenant:user", "FindByUserId")
	return ctl
}

func (s *TenantAPI) InitTenant(ctx context.Context, tenant *command.InitTenant) error {

	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {

		ctx, err := restapp.NewTestContext(context.Background(), func(option *restapp.ContextOption) {
			tenantId := service.SystemTenantId
			option.TenantId = &tenantId
		})

		if err != nil {
			return err
		}

		err = s.tenantService.InitSysTenant(ctx)
		if err != nil {
			return err
		}

		err = s.tenantService.InitDemoTenant(ctx, tenant)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *TenantAPI) Create(ctx context.Context, cmd *command.TenantCreateCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {

		err := s.tenantService.Create(ctx, &cmd.Data)
		if err != nil {
			return err
		}

		au, ok := appctx.GetAuthUser(ctx)
		if !ok {
			return errors.New("没有找到用户信息")
		}

		tu, _ := model.NewTenantUser()
		tu.Id = idutils.NewId()
		tu.TenId = cmd.Data.Id
		tu.UserId = au.GetId()
		tu.IsAdmin = true

		tus := []*model.TenantUser{tu}

		err = s.tenantUserService.CreateMany(ctx, tus)
		if err != nil {
			return err
		}

		//初始化币种
		err = s.currencyService.InitCurrency(ctx, cmd.Data.Id)
		if err != nil {
			return err
		}
		//初始化银行
		err = s.bankService.InitBank(ctx, cmd.Data.Id)
		if err != nil {
			return err
		}
		//初始化标签
		err = s.tagService.InitTag(ctx, cmd.Data.Id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *TenantAPI) Update(ctx context.Context, cmd *command.TenantUpdateCommand) error {
	return s.tenantService.Update(ctx, cmd)
}

func (s *TenantAPI) Delete(ctx context.Context, cmd *command.TenantDeleteCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		err := s.tenantUserService.DeleteByTenantIds(ctx, []string{cmd.Data.Id})
		if err != nil {
			return err
		}
		return s.tenantService.Delete(ctx, cmd)
	})
}

func (s *TenantAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Tenant, error) {
	return s.tenantService.FindById(ctx, qry.Id)
}

func (s *TenantAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.Tenant], error) {
	ctx, _ = restapp.NewTestContext(context.Background(), func(option *restapp.ContextOption) {
		tenantId := service.SystemTenantId
		option.TenantId = &tenantId
	})
	return s.tenantService.FindPaging(ctx, qry)
}

func (s *TenantAPI) FindByUserId(ctx context.Context, qry *query.FindByUserIdQuery) ([]*model.Tenant, error) {
	return s.tenantService.FindByUserId(ctx, qry.UserId)
}
