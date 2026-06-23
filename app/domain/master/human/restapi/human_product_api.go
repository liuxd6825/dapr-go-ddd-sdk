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

type HumanProductAPI struct {
	service  *service.HumanProductService
	rootPath string
}

func NewHumanProductAPI(rootPath string) *HumanProductAPI {
	return &HumanProductAPI{
		rootPath: rootPath,
		service:  service.NewHumanProductService(),
	}
}

func (s *HumanProductAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanProductAPI", s)
	controller.Post("/human-product:submitBatch", "SubmitMany")
	controller.Post("/human-product:createBatch", "CreateMany")
	controller.Put("/human-product:updateBatch", "UpdateMany")
	controller.Delete("/human-product/{id}", "DeleteById")
	controller.GetPaging("/human-product", "FindPaging")
	controller.GetData("/human/{humanId}/human-product", "FindByHumanId")
	controller.GetOne("/human-product/{id}", "FindById")
	return controller
}

func (s *HumanProductAPI) SubmitMany(ctx context.Context, cmd *command.HumanProductSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanProductAPI) CreateMany(ctx context.Context, cmd *command.HumanProductCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanProductAPI) UpdateMany(ctx context.Context, cmd *command.HumanProductUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanProductAPI) DeleteById(ctx context.Context, qry *query.HumanProductFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanProductAPI) FindById(ctx context.Context, qry *query.HumanProductFindByIdQuery) (*model.HumanProduct, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanProductAPI) FindPaging(ctx context.Context, qry *query.HumanProductFindPagingQuery) (store.FindPagingResult[*model.HumanProduct], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanProductAPI) FindByHumanId(ctx context.Context, qry *query.HumanProductFindByHumanIdQuery) (store.FindPagingResult[*model.HumanProduct], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
