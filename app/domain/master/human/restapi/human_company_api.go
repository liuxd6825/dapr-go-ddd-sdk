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

type HumanCompanyAPI struct {
	service  *service.HumanCompanyService
	rootPath string
}

func NewHumanCompanyAPI(rootPath string) *HumanCompanyAPI {
	return &HumanCompanyAPI{
		rootPath: rootPath,
		service:  service.NewHumanCompanyService(),
	}
}

func (s *HumanCompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanCompanyAPI", s)
	controller.Post("/human-company:submitBatch", "SubmitMany")
	controller.Post("/human-company:createBatch", "CreateMany")
	controller.Put("/human-company:updateBatch", "UpdateMany")
	controller.Delete("/human-company/{id}", "DeleteById")
	controller.GetPaging("/human-company", "FindPaging")
	controller.GetData("/human/{humanId}/human-company", "FindByHumanId")
	controller.GetOne("/human-company/{id}", "FindById")
	return controller
}

func (s *HumanCompanyAPI) SubmitMany(ctx context.Context, cmd *command.HumanCompanySubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanCompanyAPI) CreateMany(ctx context.Context, cmd *command.HumanCompanyCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCompanyAPI) UpdateMany(ctx context.Context, cmd *command.HumanCompanyUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCompanyAPI) DeleteById(ctx context.Context, qry *query.HumanCompanyFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanCompanyAPI) FindById(ctx context.Context, qry *query.HumanCompanyFindByIdQuery) (*model.HumanCompany, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanCompanyAPI) FindPaging(ctx context.Context, qry *query.HumanCompanyFindPagingQuery) (store.FindPagingResult[*model.HumanCompany], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanCompanyAPI) FindByHumanId(ctx context.Context, qry *query.HumanCompanyFindByHumanIdQuery) (store.FindPagingResult[*model.HumanCompany], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
