package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CardAPI struct {
	env         *env.Env
	cardService *service.CardService
	rootPath    string
}

func NewCardAPI(env *env.Env, rootPath string) *CardAPI {
	return &CardAPI{
		env:         env,
		cardService: service.NewCardService(),
		rootPath:    rootPath,
	}
}

func (s *CardAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.cardService = service.NewCardService()
	ctl := restapi.NewController(app, s.rootPath+"/card", "sys.CardAPI", s)
	ctl.Post("/card", "Create")
	ctl.Put("/card", "Update")
	ctl.Delete("/card", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/card/{id}", "FindById")
	ctl.GetPaging("/card", "FindPaging")
	return ctl
}

func (s *CardAPI) Create(ctx context.Context, cmd *command.CardCreateCommand) error {
	return s.cardService.Create(ctx, cmd)
}

func (s *CardAPI) Update(ctx context.Context, cmd *command.CardUpdateCommand) error {
	return s.cardService.Update(ctx, cmd)
}

func (s *CardAPI) Delete(ctx context.Context, cmd *command.CardDeleteCommand) error {
	return s.cardService.Delete(ctx, cmd)
}

func (s *CardAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Card, error) {
	return s.cardService.FindById(ctx, qry)
}

func (s *CardAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Card], error) {
	return s.cardService.FindPaging(ctx, qry.CaseId, qry)
}
