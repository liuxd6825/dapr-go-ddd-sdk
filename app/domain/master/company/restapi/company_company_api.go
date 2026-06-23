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

type CompanyCompanyAPI struct {
	service  *service.CompanyCompanyService
	rootPath string
}

func NewCompanyCompanyAPI(rootPath string) *CompanyCompanyAPI {
	return &CompanyCompanyAPI{
		rootPath: rootPath,
		service:  service.NewCompanyCompanyService(),
	}
}

func (s *CompanyCompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyCompanyAPI", s)
	controller.Post("/company-company:submitBatch", "SubmitMany")
	controller.Post("/company-company:createBatch", "CreateMany")
	controller.Put("/company-company:updateBatch", "UpdateMany")
	controller.Delete("/company-company/{id}", "DeleteById")
	controller.GetPaging("/company-company", "FindPaging")
	controller.GetData("/company/{companyId}/company-company", "FindByCompanyId")
	controller.GetOne("/company-company/{id}", "FindById")
	return controller
}

func (s *CompanyCompanyAPI) SubmitMany(ctx context.Context, cmd *command.CompanyCompanySubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *CompanyCompanyAPI) CreateMany(ctx context.Context, cmd *command.CompanyCompanyCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyCompanyAPI) UpdateMany(ctx context.Context, cmd *command.CompanyCompanyUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyCompanyAPI) DeleteById(ctx context.Context, qry *query.CompanyCompanyFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *CompanyCompanyAPI) FindById(ctx context.Context, qry *query.CompanyCompanyFindByIdQuery) (*model.CompanyCompany, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyCompanyAPI) FindPaging(ctx context.Context, qry *query.CompanyCompanyFindPagingQuery) (store.FindPagingResult[*model.CompanyCompany], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyCompanyAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyCompanyFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyCompany], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
