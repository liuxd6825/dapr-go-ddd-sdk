package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type FunAPI struct {
	env             *envs.Env
	funService      *service.FunService
	cardFileService *service.CardFileService
	rootPath        string
	mongoDbKey      string
}

func NewFunAPI(env *envs.Env, rootPath string, mongoDbKey string) *FunAPI {
	return &FunAPI{
		env:             env,
		funService:      service.NewFunService(mongoDbKey),
		cardFileService: service.NewCardFileService(mongoDbKey),
		rootPath:        rootPath,
	}
}

func (s *FunAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.funService = service.NewFunService(s.mongoDbKey)
	ctl := restapi.NewController(app, s.rootPath+"/card", "sys.NewFunAPI", s)
	ctl.Post("/fun", "Create")
	ctl.Put("/fun", "Update")
	ctl.Delete("/fun", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/fun/{id}", "FindById")
	ctl.GetPaging("/fun", "FindPaging")
	ctl.GetData("/fun:view", "FindViewByAppId")
	return ctl
}

func (s *FunAPI) Create(ctx context.Context, cmd *command.FunCreateCommand) error {
	return s.funService.Create(ctx, cmd)
}

func (s *FunAPI) Update(ctx context.Context, cmd *command.FunUpdateCommand) error {
	return s.funService.Update(ctx, cmd)
}

func (s *FunAPI) Delete(ctx context.Context, cmd *command.FunDeleteCommand) error {
	return s.funService.Delete(ctx, cmd)
}

func (s *FunAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Fun, error) {
	return s.funService.FindById(ctx, qry)
}

func (s *FunAPI) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Fun], error) {
	return s.funService.FindPaging(ctx, qry)
}

func (s *FunAPI) FindViewByAppId(ctx context.Context, qry *query.FindByAppIdQuery) ([]*model.FunView, error) {
	fvs, err := s.funService.FindViewByAppId(ctx, qry.AppId)
	if err != nil {
		return nil, err
	}
	for _, v := range fvs {
		cfs, err := s.cardFileService.FindByFunId(ctx, v.Id)
		if err != nil {
			return nil, err
		}
		v.Cards = cfs
	}
	return fvs, nil
}
