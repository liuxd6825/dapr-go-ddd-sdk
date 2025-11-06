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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CardFileAPI struct {
	env             *env.Env
	cardFileService *service.CardFileService
	rootPath        string
}

func NewCardFileAPI(env *env.Env, rootPath string) *CardFileAPI {
	return &CardFileAPI{
		env:             env,
		cardFileService: service.NewCardFileService(),
		rootPath:        rootPath,
	}
}

func (s *CardFileAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.cardFileService = service.NewCardFileService()
	ctl := restapi.NewController(app, s.rootPath+"/card", "sys.CardFileAPI", s)
	ctl.Post("/card-file", "Create")
	ctl.Put("/card-file", "Update")
	ctl.Delete("/card-file", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/card-file/{id}", "FindById")
	ctl.GetPaging("/card-file", "FindPaging")
	ctl.GetData("/card-file:app-fun", "FindCards")
	return ctl
}

func (s *CardFileAPI) Create(ctx context.Context, cmd *command.CardFileCreateCommand) error {
	return s.cardFileService.Create(ctx, cmd)
}

func (s *CardFileAPI) Update(ctx context.Context, cmd *command.CardFileUpdateCommand) error {
	return s.cardFileService.Update(ctx, cmd)
}

func (s *CardFileAPI) Delete(ctx context.Context, cmd *command.CardFileDeleteCommand) error {
	return s.cardFileService.Delete(ctx, cmd)
}

func (s *CardFileAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.CardFile, error) {
	return s.cardFileService.FindById(ctx, qry)
}

func (s *CardFileAPI) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.CardFile], error) {
	return s.cardFileService.FindPaging(ctx, qry)
}

func (s *CardFileAPI) FindCards(ctx context.Context, qry *query.FindCardFilesQuery) ([]*model.CardFile, error) {
	return s.cardFileService.FindCards(ctx, qry)
}
