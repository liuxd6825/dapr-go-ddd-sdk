package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type ProductHumanAPI struct {
	service  *service.ProductHumanService
	rootPath string
}

func NewProductHumanAPI(rootPath string) *ProductHumanAPI {
	return &ProductHumanAPI{
		rootPath: rootPath,
		service:  service.NewProductHumanService(),
	}
}

func (s *ProductHumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductHumanAPI", s)
	controller.Post("/product-human:submitBatch", "SubmitMany")
	controller.Post("/product-human:createBatch", "CreateMany")
	controller.Put("/product-human:updateBatch", "UpdateMany")
	controller.Delete("/product-human/{id}", "DeleteById")
	controller.GetPaging("/product-human", "FindPaging")
	controller.GetData("/product/{productId}/product-human", "FindByProductId")
	controller.GetOne("/product-human/{id}", "FindById")
	return controller
}

func (s *ProductHumanAPI) SubmitMany(ctx context.Context, cmd *command.ProductHumanSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ProductHumanAPI) CreateMany(ctx context.Context, cmd *command.ProductHumanCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductHumanAPI) UpdateMany(ctx context.Context, cmd *command.ProductHumanUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductHumanAPI) DeleteById(ctx context.Context, qry *query.ProductHumanFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ProductHumanAPI) FindById(ctx context.Context, qry *query.ProductHumanFindByIdQuery) (*model.ProductHuman, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ProductHumanAPI) FindPaging(ctx context.Context, qry *query.ProductHumanFindPagingQuery) (store.FindPagingResult[*model.ProductHuman], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ProductHumanAPI) FindByProductId(ctx context.Context, qry *query.ProductHumanFindByProductIdQuery) (store.FindPagingResult[*model.ProductHuman], error) {
	res := s.service.FindByProductId(ctx, qry, qry.ProductId)
	return res, res.GetError()
}
