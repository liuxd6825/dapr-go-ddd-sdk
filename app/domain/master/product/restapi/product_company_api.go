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

type ProductCompanyAPI struct {
	service  *service.ProductCompanyService
	rootPath string
}

func NewProductCompanyAPI(rootPath string) *ProductCompanyAPI {
	return &ProductCompanyAPI{
		rootPath: rootPath,
		service:  service.NewProductCompanyService(),
	}
}

func (s *ProductCompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductCompanyAPI", s)
	controller.Post("/product-company:submitBatch", "SubmitMany")
	controller.Post("/product-company:createBatch", "CreateMany")
	controller.Put("/product-company:updateBatch", "UpdateMany")
	controller.Delete("/product-company/{id}", "DeleteById")
	controller.GetPaging("/product-company", "FindPaging")
	controller.GetData("/product/{productId}/product-company", "FindByProductId")
	controller.GetOne("/product-company/{id}", "FindById")
	return controller
}

func (s *ProductCompanyAPI) SubmitMany(ctx context.Context, cmd *command.ProductCompanySubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ProductCompanyAPI) CreateMany(ctx context.Context, cmd *command.ProductCompanyCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductCompanyAPI) UpdateMany(ctx context.Context, cmd *command.ProductCompanyUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductCompanyAPI) DeleteById(ctx context.Context, qry *query.ProductCompanyFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ProductCompanyAPI) FindById(ctx context.Context, qry *query.ProductCompanyFindByIdQuery) (*model.ProductCompany, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ProductCompanyAPI) FindPaging(ctx context.Context, qry *query.ProductCompanyFindPagingQuery) (store.FindPagingResult[*model.ProductCompany], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ProductCompanyAPI) FindByProductId(ctx context.Context, qry *query.ProductCompanyFindByProductIdQuery) (store.FindPagingResult[*model.ProductCompany], error) {
	res := s.service.FindByProductId(ctx, qry, qry.ProductId)
	return res, res.GetError()
}
