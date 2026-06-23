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

type ContractContractAPI struct {
	service  *service.ContractContractService
	rootPath string
}

func NewContractContractAPI(rootPath string) *ContractContractAPI {
	return &ContractContractAPI{
		rootPath: rootPath,
		service:  service.NewContractContractService(),
	}
}

func (s *ContractContractAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractContractAPI", s)
	controller.Post("/contract-contract:submitBatch", "SubmitMany")
	controller.Post("/contract-contract:createBatch", "CreateMany")
	controller.Put("/contract-contract:updateBatch", "UpdateMany")
	controller.Delete("/contract-contract/{id}", "DeleteById")
	controller.GetPaging("/contract-contract", "FindPaging")
	controller.GetData("/contract/{contractId}/contract-contract", "FindByContractId")
	controller.GetOne("/contract-contract/{id}", "FindById")
	return controller
}

func (s *ContractContractAPI) SubmitMany(ctx context.Context, cmd *command.ContractContractSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ContractContractAPI) CreateMany(ctx context.Context, cmd *command.ContractContractCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractContractAPI) UpdateMany(ctx context.Context, cmd *command.ContractContractUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractContractAPI) DeleteById(ctx context.Context, qry *query.ContractContractFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ContractContractAPI) FindById(ctx context.Context, qry *query.ContractContractFindByIdQuery) (*model.ContractContract, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ContractContractAPI) FindPaging(ctx context.Context, qry *query.ContractContractFindPagingQuery) (store.FindPagingResult[*model.ContractContract], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ContractContractAPI) FindByContractId(ctx context.Context, qry *query.ContractContractFindByContractIdQuery) (store.FindPagingResult[*model.ContractContract], error) {
	res := s.service.FindByContractId(ctx, qry, qry.ContractId)
	return res, res.GetError()
}
