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

type HumanRecordAPI struct {
	service  *service.HumanRecordService
	rootPath string
}

func NewHumanRecordAPI(rootPath string) *HumanRecordAPI {
	return &HumanRecordAPI{
		rootPath: rootPath,
		service:  service.NewHumanRecordService(),
	}
}

func (s *HumanRecordAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanRecordAPI", s)
	controller.Post("/human-record:submitBatch", "SubmitMany")
	controller.Post("/human-record:createBatch", "CreateMany")
	controller.Put("/human-record:updateBatch", "UpdateMany")
	controller.Delete("/human-record/{id}", "DeleteById")
	controller.GetPaging("/human-record", "FindPaging")
	controller.GetData("/human/{humanId}/human-record", "FindByHumanId")
	controller.GetOne("/human-record/{id}", "FindById")
	return controller
}

func (s *HumanRecordAPI) SubmitMany(ctx context.Context, cmd *command.HumanRecordSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanRecordAPI) CreateMany(ctx context.Context, cmd *command.HumanRecordCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanRecordAPI) UpdateMany(ctx context.Context, cmd *command.HumanRecordUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanRecordAPI) DeleteById(ctx context.Context, qry *query.HumanRecordFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanRecordAPI) FindById(ctx context.Context, qry *query.HumanRecordFindByIdQuery) (*model.HumanRecord, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanRecordAPI) FindPaging(ctx context.Context, qry *query.HumanRecordFindPagingQuery) (store.FindPagingResult[*model.HumanRecord], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanRecordAPI) FindByHumanId(ctx context.Context, qry *query.HumanRecordFindByHumanIdQuery) (store.FindPagingResult[*model.HumanRecord], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
