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

type CompanyAccountAPI struct {
	service  *service.CompanyAccountService
	rootPath string
}

func NewCompanyAccountAPI(rootPath string) *CompanyAccountAPI {
	return &CompanyAccountAPI{
		rootPath: rootPath,
		service:  service.NewCompanyAccountService(),
	}
}

func (s *CompanyAccountAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyAccountAPI", s)
	controller.Post("/company-account:submitBatch", "SubmitMany")
	controller.Post("/company-account:createBatch", "CreateMany")
	controller.Put("/company-account:updateBatch", "UpdateMany")
	controller.Delete("/company-account/{id}", "DeleteById")
	controller.GetPaging("/company-account", "FindPaging")
	controller.GetData("/company/{companyId}/company-account", "FindByCompanyId")
	controller.GetOne("/company-account/{id}", "FindById")
	return controller
}

func (s *CompanyAccountAPI) SubmitMany(ctx context.Context, cmd *command.CompanyAccountSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *CompanyAccountAPI) CreateMany(ctx context.Context, cmd *command.CompanyAccountCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyAccountAPI) UpdateMany(ctx context.Context, cmd *command.CompanyAccountUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyAccountAPI) DeleteById(ctx context.Context, qry *query.CompanyAccountFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *CompanyAccountAPI) FindById(ctx context.Context, qry *query.CompanyAccountFindByIdQuery) (*model.CompanyAccount, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyAccountAPI) FindPaging(ctx context.Context, qry *query.CompanyAccountFindPagingQuery) (store.FindPagingResult[*model.CompanyAccount], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyAccountAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyAccountFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyAccount], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
