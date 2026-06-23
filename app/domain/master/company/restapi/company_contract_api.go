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

type CompanyContractAPI struct {
	service  *service.CompanyContractService
	rootPath string
}

func NewCompanyContractAPI(rootPath string) *CompanyContractAPI {
	return &CompanyContractAPI{
		rootPath: rootPath,
		service:  service.NewCompanyContractService(),
	}
}

func (s *CompanyContractAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyContractAPI", s)
	controller.Post("/company-contract:submitBatch", "SubmitMany")
	controller.Post("/company-contract:createBatch", "CreateMany")
	controller.Put("/company-contract:updateBatch", "UpdateMany")
	controller.Delete("/company-contract/{id}", "DeleteById")
	controller.GetPaging("/company-contract", "FindPaging")
	controller.GetData("/company/{companyId}/company-contract", "FindByCompanyId")
	controller.GetOne("/company-contract/{id}", "FindById")
	return controller
}

func (s *CompanyContractAPI) SubmitMany(ctx context.Context, cmd *command.CompanyContractSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *CompanyContractAPI) CreateMany(ctx context.Context, cmd *command.CompanyContractCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyContractAPI) UpdateMany(ctx context.Context, cmd *command.CompanyContractUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyContractAPI) DeleteById(ctx context.Context, qry *query.CompanyContractFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *CompanyContractAPI) FindById(ctx context.Context, qry *query.CompanyContractFindByIdQuery) (*model.CompanyContract, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyContractAPI) FindPaging(ctx context.Context, qry *query.CompanyContractFindPagingQuery) (store.FindPagingResult[*model.CompanyContract], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyContractAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyContractFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyContract], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
