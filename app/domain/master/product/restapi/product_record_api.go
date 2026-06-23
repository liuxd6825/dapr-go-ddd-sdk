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

type ProductRecordAPI struct {
	service  *service.ProductRecordService
	rootPath string
}

func NewProductRecordAPI(rootPath string) *ProductRecordAPI {
	return &ProductRecordAPI{
		rootPath: rootPath,
		service:  service.NewProductRecordService(),
	}
}

func (s *ProductRecordAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductRecordAPI", s)
	controller.Post("/product-record:submitBatch", "SubmitMany")
	controller.Post("/product-record:createBatch", "CreateMany")
	controller.Put("/product-record:updateBatch", "UpdateMany")
	controller.Delete("/product-record/{id}", "DeleteById")
	controller.GetPaging("/product-record", "FindPaging")
	controller.GetData("/product/{productId}/product-record", "FindByProductId")
	controller.GetOne("/product-record/{id}", "FindById")
	return controller
}

func (s *ProductRecordAPI) SubmitMany(ctx context.Context, cmd *command.ProductRecordSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ProductRecordAPI) CreateMany(ctx context.Context, cmd *command.ProductRecordCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductRecordAPI) UpdateMany(ctx context.Context, cmd *command.ProductRecordUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ProductRecordAPI) DeleteById(ctx context.Context, qry *query.ProductRecordFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ProductRecordAPI) FindById(ctx context.Context, qry *query.ProductRecordFindByIdQuery) (*model.ProductRecord, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ProductRecordAPI) FindPaging(ctx context.Context, qry *query.ProductRecordFindPagingQuery) (store.FindPagingResult[*model.ProductRecord], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ProductRecordAPI) FindByProductId(ctx context.Context, qry *query.ProductRecordFindByProductIdQuery) (store.FindPagingResult[*model.ProductRecord], error) {
	res := s.service.FindByProductId(ctx, qry, qry.ProductId)
	return res, res.GetError()
}
