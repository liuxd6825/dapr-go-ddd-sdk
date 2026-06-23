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

type ProductProductAPI struct {
	service  *service.ProductProductService
	rootPath string
}

func NewProductProductAPI(rootPath string) *ProductProductAPI {
	return &ProductProductAPI{
		rootPath: rootPath,
		service:  service.NewProductProductService(),
	}
}

func (s *ProductProductAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductProductAPI", s)
	controller.Post("/product-product:submitBatch", "SubmitMany")
	controller.Post("/product-product:createBatch", "CreateMany")
	controller.Put("/product-product:updateBatch", "UpdateMany")
	controller.Delete("/product-product/{id}", "DeleteById")
	controller.GetPaging("/product-product", "FindPaging")
	controller.GetData("/product/{productId}/product-product", "FindByProductId")
	controller.GetOne("/product-product/{id}", "FindById")
	return controller
}

func (s *ProductProductAPI) SubmitMany(ctx context.Context, cmd *command.ProductProductSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ProductProductAPI) CreateMany(ctx context.Context, cmd *command.ProductProductCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductProductAPI) UpdateMany(ctx context.Context, cmd *command.ProductProductUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductProductAPI) DeleteById(ctx context.Context, qry *query.ProductProductFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ProductProductAPI) FindById(ctx context.Context, qry *query.ProductProductFindByIdQuery) (*model.ProductProduct, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ProductProductAPI) FindPaging(ctx context.Context, qry *query.ProductProductFindPagingQuery) (store.FindPagingResult[*model.ProductProduct], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ProductProductAPI) FindByProductId(ctx context.Context, qry *query.ProductProductFindByProductIdQuery) (store.FindPagingResult[*model.ProductProduct], error) {
	res := s.service.FindByProductId(ctx, qry, qry.ProductId)
	return res, res.GetError()
}
