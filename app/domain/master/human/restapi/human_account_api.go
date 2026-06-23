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

type HumanAccountAPI struct {
	service  *service.HumanAccountService
	rootPath string
}

func NewHumanAccountAPI(rootPath string) *HumanAccountAPI {
	return &HumanAccountAPI{
		rootPath: rootPath,
		service:  service.NewHumanAccountService(),
	}
}

func (s *HumanAccountAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanAccountAPI", s)
	controller.Post("/human-account:submitBatch", "SubmitMany")
	controller.Post("/human-account:createBatch", "CreateMany")
	controller.Put("/human-account:updateBatch", "UpdateMany")
	controller.Delete("/human-account/{id}", "DeleteById")
	controller.GetPaging("/human-account", "FindPaging")
	controller.GetData("/human/{humanId}/human-account", "FindByHumanId")
	controller.GetOne("/human-account/{id}", "FindById")
	return controller
}

func (s *HumanAccountAPI) SubmitMany(ctx context.Context, cmd *command.HumanAccountSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanAccountAPI) CreateMany(ctx context.Context, cmd *command.HumanAccountCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanAccountAPI) UpdateMany(ctx context.Context, cmd *command.HumanAccountUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanAccountAPI) DeleteById(ctx context.Context, qry *query.HumanAccountFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanAccountAPI) FindById(ctx context.Context, qry *query.HumanAccountFindByIdQuery) (*model.HumanAccount, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanAccountAPI) FindPaging(ctx context.Context, qry *query.HumanAccountFindPagingQuery) (store.FindPagingResult[*model.HumanAccount], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanAccountAPI) FindByHumanId(ctx context.Context, qry *query.HumanAccountFindByHumanIdQuery) (store.FindPagingResult[*model.HumanAccount], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
