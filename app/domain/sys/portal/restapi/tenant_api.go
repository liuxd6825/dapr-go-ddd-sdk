package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TenantAPI struct {
	env           *env.Env
	tenantService *service.TenantService
	rootPath      string
}

func NewTenantAPI(env *env.Env, rootPath string) *TenantAPI {
	return &TenantAPI{
		env:           env,
		tenantService: service.NewTenantService(),
		rootPath:      rootPath,
	}
}

func (s *TenantAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.tenantService = service.NewTenantService()
	ctl := restapi.NewController(app, s.rootPath+"/portal", "sys.TenantAPI", s)
	ctl.Post("/tenant", "Create")
	ctl.Put("/tenant", "Update")
	ctl.Delete("/tenant", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/tenant/{id}", "FindById")
	ctl.GetPaging("/tenant", "FindPaging")
	return ctl
}

func (s *TenantAPI) Create(ctx context.Context, cmd *command.TenantCreateCommand) error {
	return s.tenantService.Create(ctx, cmd)
}

func (s *TenantAPI) Update(ctx context.Context, cmd *command.TenantUpdateCommand) error {
	return s.tenantService.Update(ctx, cmd)
}

func (s *TenantAPI) Delete(ctx context.Context, cmd *command.TenantDeleteCommand) error {
	return s.tenantService.Delete(ctx, cmd)
}

func (s *TenantAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Tenant, error) {
	return s.tenantService.FindById(ctx, qry)
}

func (s *TenantAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Tenant], error) {
	return s.tenantService.FindPaging(ctx, qry.CaseId, qry)
}
