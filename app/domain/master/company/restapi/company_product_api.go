package restapi

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/service"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CompanyProductAPI struct {
	service  *service.CompanyProductService
	rootPath string
}

func NewCompanyProductAPI(rootPath string) *CompanyProductAPI {
	return &CompanyProductAPI{
		rootPath: rootPath,
		service:  service.NewCompanyProductService(),
	}
}

func (s *CompanyProductAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyProductAPI", s)
	controller.Post("/company-product:submitBatch", "SubmitMany")
	controller.Post("/company-product:createBatch", "CreateMany")
	controller.Put("/company-product:updateBatch", "UpdateMany")
	controller.Delete("/company-product/{id}", "DeleteById")
	controller.GetPaging("/company-product", "FindPaging")
	controller.GetData("/company/{companyId}/company-product", "FindByCompanyId")
	controller.GetOne("/company-product/{id}", "FindById")
	return controller
}

func (s *CompanyProductAPI) SubmitMany(ctx context.Context, cmd *command.CompanyProductSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *CompanyProductAPI) CreateMany(ctx context.Context, cmd *command.CompanyProductCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyProductAPI) UpdateMany(ctx context.Context, cmd *command.CompanyProductUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyProductAPI) DeleteById(ctx context.Context, qry *query.CompanyProductFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *CompanyProductAPI) FindById(ctx context.Context, qry *query.CompanyProductFindByIdQuery) (*model.CompanyProduct, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyProductAPI) FindPaging(ctx context.Context, qry *query.CompanyProductFindPagingQuery) (store.FindPagingResult[*model.CompanyProduct], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyProductAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyProductFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyProduct], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
