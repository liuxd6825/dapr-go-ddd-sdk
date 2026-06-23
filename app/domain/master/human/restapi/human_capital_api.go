package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type HumanCapitalAPI struct {
	service  *service.HumanCapitalService
	rootPath string
}

func NewHumanCapitalAPI(rootPath string) *HumanCapitalAPI {
	return &HumanCapitalAPI{
		rootPath: rootPath,
		service:  service.NewHumanCapitalService(),
	}
}

func (s *HumanCapitalAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanCapitalAPI", s)
	controller.Post("/human-capital:submitBatch", "SubmitMany")
	controller.Post("/human-capital:createBatch", "CreateMany")
	controller.Put("/human-capital:updateBatch", "UpdateMany")
	controller.Delete("/human-capital/{id}", "DeleteById")
	controller.GetPaging("/human-capital", "FindPaging")
	controller.GetData("/human/{humanId}/human-capital", "FindByHumanId")
	controller.GetOne("/human-capital/{id}", "FindById")
	return controller
}

func (s *HumanCapitalAPI) SubmitMany(ctx context.Context, cmd *command.HumanCapitalSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanCapitalAPI) CreateMany(ctx context.Context, cmd *command.HumanCapitalCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCapitalAPI) UpdateMany(ctx context.Context, cmd *command.HumanCapitalUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCapitalAPI) DeleteById(ctx context.Context, qry *query.HumanCapitalFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanCapitalAPI) FindById(ctx context.Context, qry *query.HumanCapitalFindByIdQuery) (*model.HumanCapital, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanCapitalAPI) FindPaging(ctx context.Context, qry *query.HumanCapitalFindPagingQuery) (store.FindPagingResult[*model.HumanCapital], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanCapitalAPI) FindByHumanId(ctx context.Context, qry *query.HumanCapitalFindByHumanIdQuery) (store.FindPagingResult[*model.HumanCapital], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
