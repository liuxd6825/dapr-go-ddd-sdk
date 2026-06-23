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

type HumanCredentialAPI struct {
	service  *service.HumanCredentialService
	rootPath string
}

func NewHumanCredentialAPI(rootPath string) *HumanCredentialAPI {
	return &HumanCredentialAPI{
		rootPath: rootPath,
		service:  service.NewHumanCredentialService(),
	}
}

func (s *HumanCredentialAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanCredentialAPI", s)
	controller.Post("/human-credential:submitBatch", "SubmitMany")
	controller.Post("/human-credential:createBatch", "CreateMany")
	controller.Put("/human-credential:updateBatch", "UpdateMany")
	controller.Delete("/human-credential/{id}", "DeleteById")
	controller.GetPaging("/human-credential", "FindPaging")
	controller.GetData("/human/{humanId}/human-credential", "FindByHumanId")
	controller.GetOne("/human-credential/{id}", "FindById")
	return controller
}

func (s *HumanCredentialAPI) SubmitMany(ctx context.Context, cmd *command.HumanCredentialSubmitManyCommand) (any, error) {
	return s.service.SubmitMany(ctx, cmd.Data.InsertData, cmd.Data.UpdateData)
}

func (s *HumanCredentialAPI) CreateMany(ctx context.Context, cmd *command.HumanCredentialCreateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.CreateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCredentialAPI) UpdateMany(ctx context.Context, cmd *command.HumanCredentialUpdateCommand) (any, error) {
	ids := make([]string, 0, len(cmd.Data))
	for _, item := range cmd.Data {
		ids = append(ids, item.Id)
	}
	if err := s.service.UpdateMany(ctx, cmd.Data); err != nil {
		return nil, err
	}
	return s.service.FindByIds(ctx, ids)
}

func (s *HumanCredentialAPI) DeleteById(ctx context.Context, qry *query.HumanCredentialFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanCredentialAPI) FindById(ctx context.Context, qry *query.HumanCredentialFindByIdQuery) (*model.HumanCredential, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanCredentialAPI) FindPaging(ctx context.Context, qry *query.HumanCredentialFindPagingQuery) (store.FindPagingResult[*model.HumanCredential], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *HumanCredentialAPI) FindByHumanId(ctx context.Context, qry *query.HumanCredentialFindByHumanIdQuery) (store.FindPagingResult[*model.HumanCredential], error) {
	res := s.service.FindByHumanId(ctx, qry, qry.HumanId)
	return res, res.GetError()
}
