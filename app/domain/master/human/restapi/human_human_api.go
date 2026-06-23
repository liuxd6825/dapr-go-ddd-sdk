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

type HumanHumanAPI struct {
	service  *service.HumanHumanService
	rootPath string
}

func NewHumanHumanAPI(rootPath string) *HumanHumanAPI {
	return &HumanHumanAPI{
		rootPath: rootPath,
		service:  service.NewHumanHumanService(),
	}
}

func (s *HumanHumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanHumanAPI", s)
	controller.Post("/human-human:submitBatch", "SubmitMany")
	controller.Post("/human-human:createBatch", "CreateMany")
	controller.Put("/human-human:updateBatch", "UpdateMany")
	controller.Delete("/human-human/{id}", "DeleteById")
	controller.GetPaging("/human-human", "FindPaging")
	controller.GetData("/human/{humanId}/human-human", "FindByHumanId")
	controller.GetOne("/human-human/{id}", "FindById")
	return controller
}

func (s *HumanHumanAPI) SubmitMany(ctx context.Context, cmd *command.HumanHumanSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanHumanAPI) CreateMany(ctx context.Context, cmd *command.HumanHumanCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanHumanAPI) UpdateMany(ctx context.Context, cmd *command.HumanHumanUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanHumanAPI) DeleteById(ctx context.Context, qry *query.HumanHumanFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanHumanAPI) FindById(ctx context.Context, qry *query.HumanHumanFindByIdQuery) (*model.HumanHuman, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanHumanAPI) FindPaging(ctx context.Context, qry *query.HumanHumanFindPagingQuery) (store.FindPagingResult[*model.HumanHuman], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanHumanAPI) FindByHumanId(ctx context.Context, qry *query.HumanHumanFindByHumanIdQuery) (store.FindPagingResult[*model.HumanHuman], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
