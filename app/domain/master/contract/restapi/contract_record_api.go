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

type ContractRecordAPI struct {
	service  *service.ContractRecordService
	rootPath string
}

func NewContractRecordAPI(rootPath string) *ContractRecordAPI {
	return &ContractRecordAPI{
		rootPath: rootPath,
		service:  service.NewContractRecordService(),
	}
}

func (s *ContractRecordAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractRecordAPI", s)
	controller.Post("/contract-record:submitBatch", "SubmitMany")
	controller.Post("/contract-record:createBatch", "CreateMany")
	controller.Put("/contract-record:updateBatch", "UpdateMany")
	controller.Delete("/contract-record/{id}", "DeleteById")
	controller.GetPaging("/contract-record", "FindPaging")
	controller.GetData("/contract/{contractId}/contract-record", "FindByContractId")
	controller.GetOne("/contract-record/{id}", "FindById")
	return controller
}

func (s *ContractRecordAPI) SubmitMany(ctx context.Context, cmd *command.ContractRecordSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *ContractRecordAPI) CreateMany(ctx context.Context, cmd *command.ContractRecordCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractRecordAPI) UpdateMany(ctx context.Context, cmd *command.ContractRecordUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *ContractRecordAPI) DeleteById(ctx context.Context, qry *query.ContractRecordFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *ContractRecordAPI) FindById(ctx context.Context, qry *query.ContractRecordFindByIdQuery) (*model.ContractRecord, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ContractRecordAPI) FindPaging(ctx context.Context, qry *query.ContractRecordFindPagingQuery) (store.FindPagingResult[*model.ContractRecord], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ContractRecordAPI) FindByContractId(ctx context.Context, qry *query.ContractRecordFindByContractIdQuery) (store.FindPagingResult[*model.ContractRecord], error) {
	res := s.service.FindByContractId(ctx, qry, qry.ContractId)
	return res, res.GetError()
}
