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

type ProductContractAPI struct {
	service  *service.ProductContractService
	rootPath string
}

func NewProductContractAPI(rootPath string) *ProductContractAPI {
	return &ProductContractAPI{
		rootPath: rootPath,
		service:  service.NewProductContractService(),
	}
}

func (s *ProductContractAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductContractAPI", s)
	controller.Post("/product-contract:submitBatch", "SubmitMany")
	controller.Post("/product-contract:createBatch", "CreateMany")
	controller.Put("/product-contract:updateBatch", "UpdateMany")
	controller.Delete("/product-contract/{id}", "DeleteById")
	controller.GetPaging("/product-contract", "FindPaging")
	controller.GetData("/product/{productId}/product-contract", "FindByProductId")
	controller.GetOne("/product-contract/{id}", "FindById")
	return controller
}

func (s *ProductContractAPI) SubmitMany(ctx context.Context, cmd *command.ProductContractSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ProductContractAPI) CreateMany(ctx context.Context, cmd *command.ProductContractCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductContractAPI) UpdateMany(ctx context.Context, cmd *command.ProductContractUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductContractAPI) DeleteById(ctx context.Context, qry *query.ProductContractFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ProductContractAPI) FindById(ctx context.Context, qry *query.ProductContractFindByIdQuery) (*model.ProductContract, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ProductContractAPI) FindPaging(ctx context.Context, qry *query.ProductContractFindPagingQuery) (store.FindPagingResult[*model.ProductContract], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ProductContractAPI) FindByProductId(ctx context.Context, qry *query.ProductContractFindByProductIdQuery) (store.FindPagingResult[*model.ProductContract], error) {
	res := s.service.FindByProductId(ctx, qry, qry.ProductId)
	return res, res.GetError()
}
