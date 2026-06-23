package restapi

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/service"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CompanyHumanAPI struct {
	service  *service.CompanyHumanService
	rootPath string
}

func NewCompanyHumanAPI(rootPath string) *CompanyHumanAPI {
	return &CompanyHumanAPI{
		rootPath: rootPath,
		service:  service.NewCompanyHumanService(),
	}
}

func (s *CompanyHumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyHumanAPI", s)
	controller.Post("/company-human:submitBatch", "SubmitMany")
	controller.Post("/company-human:createBatch", "CreateMany")
	controller.Put("/company-human:updateBatch", "UpdateMany")
	controller.Delete("/company-human/{id}", "DeleteById")
	controller.GetPaging("/company-human", "FindPaging")
	controller.GetData("/company/{companyId}/company-human", "FindByCompanyId")
	controller.GetOne("/company-human/{id}", "FindById")
	return controller
}

func (s *CompanyHumanAPI) SubmitMany(ctx context.Context, cmd *command.CompanyHumanSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *CompanyHumanAPI) CreateMany(ctx context.Context, cmd *command.CompanyHumanCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyHumanAPI) UpdateMany(ctx context.Context, cmd *command.CompanyHumanUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *CompanyHumanAPI) DeleteById(ctx context.Context, qry *query.CompanyHumanFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *CompanyHumanAPI) FindById(ctx context.Context, qry *query.CompanyHumanFindByIdQuery) (*model.CompanyHuman, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyHumanAPI) FindPaging(ctx context.Context, qry *query.CompanyHumanFindPagingQuery) (store.FindPagingResult[*model.CompanyHuman], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyHumanAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyHumanFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyHuman], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
