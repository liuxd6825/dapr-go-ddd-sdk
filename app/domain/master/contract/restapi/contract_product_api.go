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

type ContractProductAPI struct {
	service  *service.ContractProductService
	rootPath string
}

func NewContractProductAPI(rootPath string) *ContractProductAPI {
	return &ContractProductAPI{
		rootPath: rootPath,
		service:  service.NewContractProductService(),
	}
}

func (s *ContractProductAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractProductAPI", s)
	controller.Post("/contract-product:submitBatch", "SubmitMany")
	controller.Post("/contract-product:createBatch", "CreateMany")
	controller.Put("/contract-product:updateBatch", "UpdateMany")
	controller.Delete("/contract-product/{id}", "DeleteById")
	controller.GetPaging("/contract-product", "FindPaging")
	controller.GetData("/contract/{contractId}/contract-product", "FindByContractId")
	controller.GetOne("/contract-product/{id}", "FindById")
	return controller
}

func (s *ContractProductAPI) SubmitMany(ctx context.Context, cmd *command.ContractProductSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ContractProductAPI) CreateMany(ctx context.Context, cmd *command.ContractProductCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractProductAPI) UpdateMany(ctx context.Context, cmd *command.ContractProductUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractProductAPI) DeleteById(ctx context.Context, qry *query.ContractProductFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ContractProductAPI) FindById(ctx context.Context, qry *query.ContractProductFindByIdQuery) (*model.ContractProduct, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ContractProductAPI) FindPaging(ctx context.Context, qry *query.ContractProductFindPagingQuery) (store.FindPagingResult[*model.ContractProduct], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ContractProductAPI) FindByContractId(ctx context.Context, qry *query.ContractProductFindByContractIdQuery) (store.FindPagingResult[*model.ContractProduct], error) {
	res := s.service.FindByContractId(ctx, qry, qry.ContractId)
	return res, res.GetError()
}
