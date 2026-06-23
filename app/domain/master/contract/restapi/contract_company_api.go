package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type ContractCompanyAPI struct {
	service  *service.ContractCompanyService
	rootPath string
}

func NewContractCompanyAPI(rootPath string) *ContractCompanyAPI {
	return &ContractCompanyAPI{
		rootPath: rootPath,
		service:  service.NewContractCompanyService(),
	}
}

func (s *ContractCompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractCompanyAPI", s)
	controller.Post("/contract-company:submitBatch", "SubmitMany")
	controller.Post("/contract-company:createBatch", "CreateMany")
	controller.Put("/contract-company:updateBatch", "UpdateMany")
	controller.Delete("/contract-company/{id}", "DeleteById")
	controller.GetPaging("/contract-company", "FindPaging")
	controller.GetData("/contract/{contractId}/contract-company", "FindByContractId")
	controller.GetOne("/contract-company/{id}", "FindById")
	return controller
}

func (s *ContractCompanyAPI) SubmitMany(ctx context.Context, cmd *command.ContractCompanySubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ContractCompanyAPI) CreateMany(ctx context.Context, cmd *command.ContractCompanyCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractCompanyAPI) UpdateMany(ctx context.Context, cmd *command.ContractCompanyUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractCompanyAPI) DeleteById(ctx context.Context, qry *query.ContractCompanyFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ContractCompanyAPI) FindById(ctx context.Context, qry *query.ContractCompanyFindByIdQuery) (*model.ContractCompany, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ContractCompanyAPI) FindPaging(ctx context.Context, qry *query.ContractCompanyFindPagingQuery) (store.FindPagingResult[*model.ContractCompany], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ContractCompanyAPI) FindByContractId(ctx context.Context, qry *query.ContractCompanyFindByContractIdQuery) (store.FindPagingResult[*model.ContractCompany], error) {
	res := s.service.FindByContractId(ctx, qry, qry.ContractId)
	return res, res.GetError()
}
