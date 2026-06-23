package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type HumanContractAPI struct {
	service  *service.HumanContractService
	rootPath string
}

func NewHumanContractAPI(rootPath string) *HumanContractAPI {
	return &HumanContractAPI{
		rootPath: rootPath,
		service:  service.NewHumanContractService(),
	}
}

func (s *HumanContractAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanContractAPI", s)
	controller.Post("/human-contract:submitBatch", "SubmitMany")
	controller.Post("/human-contract:createBatch", "CreateMany")
	controller.Put("/human-contract:updateBatch", "UpdateMany")
	controller.Delete("/human-contract/{id}", "DeleteById")
	controller.GetPaging("/human-contract", "FindPaging")
	controller.GetData("/human/{humanId}/human-contract", "FindByHumanId")
	controller.GetOne("/human-contract/{id}", "FindById")
	return controller
}

func (s *HumanContractAPI) SubmitMany(ctx context.Context, cmd *command.HumanContractSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanContractAPI) CreateMany(ctx context.Context, cmd *command.HumanContractCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanContractAPI) UpdateMany(ctx context.Context, cmd *command.HumanContractUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanContractAPI) DeleteById(ctx context.Context, qry *query.HumanContractFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanContractAPI) FindById(ctx context.Context, qry *query.HumanContractFindByIdQuery) (*model.HumanContract, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanContractAPI) FindPaging(ctx context.Context, qry *query.HumanContractFindPagingQuery) (store.FindPagingResult[*model.HumanContract], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanContractAPI) FindByHumanId(ctx context.Context, qry *query.HumanContractFindByHumanIdQuery) (store.FindPagingResult[*model.HumanContract], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
