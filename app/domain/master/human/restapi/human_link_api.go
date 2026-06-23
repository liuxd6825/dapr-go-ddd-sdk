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

type HumanLinkAPI struct {
	service  *service.HumanLinkService
	rootPath string
}

func NewHumanLinkAPI(rootPath string) *HumanLinkAPI {
	return &HumanLinkAPI{
		rootPath: rootPath,
		service:  service.NewHumanLinkService(),
	}
}

func (s *HumanLinkAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanLinkAPI", s)
	controller.Post("/human-link:submitBatch", "SubmitMany")
	controller.Post("/human-link:createBatch", "CreateMany")
	controller.Put("/human-link:updateBatch", "UpdateMany")
	controller.Delete("/human-link/{id}", "DeleteById")
	controller.GetPaging("/human-link", "FindPaging")
	controller.GetData("/human/{humanId}/human-link", "FindByHumanId")
	controller.GetOne("/human-link/{id}", "FindById")
	return controller
}

func (s *HumanLinkAPI) SubmitMany(ctx context.Context, cmd *command.HumanLinkSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanLinkAPI) CreateMany(ctx context.Context, cmd *command.HumanLinkCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanLinkAPI) UpdateMany(ctx context.Context, cmd *command.HumanLinkUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanLinkAPI) DeleteById(ctx context.Context, qry *query.HumanLinkFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanLinkAPI) FindById(ctx context.Context, qry *query.HumanLinkFindByIdQuery) (*model.HumanLink, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanLinkAPI) FindPaging(ctx context.Context, qry *query.HumanLinkFindPagingQuery) (store.FindPagingResult[*model.HumanLink], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanLinkAPI) FindByHumanId(ctx context.Context, qry *query.HumanLinkFindByHumanIdQuery) (store.FindPagingResult[*model.HumanLink], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
