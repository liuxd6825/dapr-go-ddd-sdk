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

type HumanAddressAPI struct {
	service  *service.HumanAddressService
	rootPath string
}

func NewHumanAddressAPI(rootPath string) *HumanAddressAPI {
	return &HumanAddressAPI{
		rootPath: rootPath,
		service:  service.NewHumanAddressService(),
	}
}

func (s *HumanAddressAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanAddressAPI", s)
	controller.Post("/human-address:submitBatch", "SubmitMany")
	controller.Post("/human-address:createBatch", "CreateMany")
	controller.Put("/human-address:updateBatch", "UpdateMany")
	controller.Delete("/human-address/{id}", "DeleteById")
	controller.GetPaging("/human-address", "FindPaging")
	controller.GetData("/human/{humanId}/human-address", "FindByHumanId")
	controller.GetOne("/human-address/{id}", "FindById")
	return controller
}

func (s *HumanAddressAPI) SubmitMany(ctx context.Context, cmd *command.HumanAddressSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanAddressAPI) CreateMany(ctx context.Context, cmd *command.HumanAddressCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanAddressAPI) UpdateMany(ctx context.Context, cmd *command.HumanAddressUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanAddressAPI) DeleteById(ctx context.Context, qry *query.HumanAddressFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanAddressAPI) FindById(ctx context.Context, qry *query.HumanAddressFindByIdQuery) (*model.HumanAddress, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanAddressAPI) FindPaging(ctx context.Context, qry *query.HumanAddressFindPagingQuery) (store.FindPagingResult[*model.HumanAddress], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanAddressAPI) FindByHumanId(ctx context.Context, qry *query.HumanAddressFindByHumanIdQuery) (store.FindPagingResult[*model.HumanAddress], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
