package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type AppAPI struct {
	env        *env.Env
	appService *service.AppService
	funService *service.FunService
	rootPath   string
}

func NewAppAPI(env *env.Env, rootPath string) *AppAPI {
	return &AppAPI{
		env:        env,
		appService: service.NewAppService(),
		funService: service.NewFunService(),
		rootPath:   rootPath,
	}
}

func (s *AppAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.appService = service.NewAppService()
	ctl := restapi.NewController(app, s.rootPath+"/card", "sys.AppAPI", s)
	ctl.Post("/app", "Create")
	ctl.Put("/app", "Update")
	ctl.Delete("/app", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/app/{id}", "FindById")
	ctl.GetPaging("/app", "FindPaging")
	ctl.GetData("/app:tree", "FindTree")
	return ctl
}

func (s *AppAPI) Create(ctx context.Context, cmd *command.AppCreateCommand) error {
	return s.appService.Create(ctx, cmd)
}

func (s *AppAPI) Update(ctx context.Context, cmd *command.AppUpdateCommand) error {
	return s.appService.Update(ctx, cmd)
}

func (s *AppAPI) Delete(ctx context.Context, cmd *command.AppDeleteCommand) error {
	return s.appService.Delete(ctx, cmd)
}

func (s *AppAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.App, error) {
	return s.appService.FindById(ctx, qry)
}

func (s *AppAPI) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.App], error) {
	return s.appService.FindPaging(ctx, qry)
}

func (s *AppAPI) FindTree(ctx context.Context) ([]*model.AppTree, error) {
	apps, err := s.appService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	funs, err := s.funService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.appService.TransformTree(ctx, apps, funs), nil
}
