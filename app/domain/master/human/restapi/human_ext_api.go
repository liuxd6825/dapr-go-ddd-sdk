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

type HumanExtAPI struct {
	service  *service.HumanExtService
	rootPath string
}

func NewHumanExtAPI(rootPath string) *HumanExtAPI {
	return &HumanExtAPI{
		rootPath: rootPath,
		service:  service.NewHumanExtService(),
	}
}

func (s *HumanExtAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanExtAPI", s)
	controller.Post("/human-ext:submitBatch", "SubmitMany")
	controller.Post("/human-ext:createBatch", "CreateMany")
	controller.Put("/human-ext:updateBatch", "UpdateMany")
	controller.Delete("/human-ext/{id}", "DeleteById")
	controller.GetPaging("/human-ext", "FindPaging")
	controller.GetData("/human/{humanId}/human-ext", "FindByHumanId")
	controller.GetOne("/human-ext/{id}", "FindById")
	return controller
}

func (s *HumanExtAPI) SubmitMany(ctx context.Context, cmd *command.HumanExtSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanExtAPI) CreateMany(ctx context.Context, cmd *command.HumanExtCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanExtAPI) UpdateMany(ctx context.Context, cmd *command.HumanExtUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanExtAPI) DeleteById(ctx context.Context, qry *query.HumanExtFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanExtAPI) FindById(ctx context.Context, qry *query.HumanExtFindByIdQuery) (*model.HumanExt, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanExtAPI) FindPaging(ctx context.Context, qry *query.HumanExtFindPagingQuery) (store.FindPagingResult[*model.HumanExt], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanExtAPI) FindByHumanId(ctx context.Context, qry *query.HumanExtFindByHumanIdQuery) (store.FindPagingResult[*model.HumanExt], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
