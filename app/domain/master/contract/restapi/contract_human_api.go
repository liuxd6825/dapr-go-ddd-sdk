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

type ContractHumanAPI struct {
	service  *service.ContractHumanService
	rootPath string
}

func NewContractHumanAPI(rootPath string) *ContractHumanAPI {
	return &ContractHumanAPI{
		rootPath: rootPath,
		service:  service.NewContractHumanService(),
	}
}

func (s *ContractHumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractHumanAPI", s)
	controller.Post("/contract-human:submitBatch", "SubmitMany")
	controller.Post("/contract-human:createBatch", "CreateMany")
	controller.Put("/contract-human:updateBatch", "UpdateMany")
	controller.Delete("/contract-human/{id}", "DeleteById")
	controller.GetPaging("/contract-human", "FindPaging")
	controller.GetData("/contract/{contractId}/contract-human", "FindByContractId")
	controller.GetOne("/contract-human/{id}", "FindById")
	return controller
}

func (s *ContractHumanAPI) SubmitMany(ctx context.Context, cmd *command.ContractHumanSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ContractHumanAPI) CreateMany(ctx context.Context, cmd *command.ContractHumanCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractHumanAPI) UpdateMany(ctx context.Context, cmd *command.ContractHumanUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractHumanAPI) DeleteById(ctx context.Context, qry *query.ContractHumanFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ContractHumanAPI) FindById(ctx context.Context, qry *query.ContractHumanFindByIdQuery) (*model.ContractHuman, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ContractHumanAPI) FindPaging(ctx context.Context, qry *query.ContractHumanFindPagingQuery) (store.FindPagingResult[*model.ContractHuman], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ContractHumanAPI) FindByContractId(ctx context.Context, qry *query.ContractHumanFindByContractIdQuery) (store.FindPagingResult[*model.ContractHuman], error) {
	res := s.service.FindByContractId(ctx, qry, qry.ContractId)
	return res, res.GetError()
}
