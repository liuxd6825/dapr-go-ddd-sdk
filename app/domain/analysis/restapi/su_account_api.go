package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/command"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type SuTaskAccountApi struct {
	accountService *service.SuAccountService
	env            *env.Env
	rootPath       string
}

func NewSuTaskAccountApi(env *env.Env, rootPath string) *SuTaskAccountApi {
	accountService := service.NewSuTaskAccountService()
	return &SuTaskAccountApi{
		env:            env,
		rootPath:       rootPath,
		accountService: accountService,
	}
}

func (s *SuTaskAccountApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.SuAccountApi", s)
	controller.Post("su-account", "Create")
	controller.Post("su-account:batch", "BatchCreate")
	controller.Put("su-account", "Update")
	controller.Delete("su-account:batch", "BatchDelete", restapi.WithParamsInBody(true))
	controller.GetOne("su-account/{id}", "FindById")
	controller.GetPaging("su-account", "FindPaging")
	return controller
}

func (s *SuTaskAccountApi) Create(ctx context.Context, cmd *command.SuTaskAccountCreateCommand) error {
	return s.accountService.Create(ctx, &cmd.Data)
}

func (s *SuTaskAccountApi) BatchCreate(ctx context.Context, cmd *command.SuTaskAccountBatchCreateCommand) error {
	return s.accountService.CreateMany(ctx, cmd.Data)
}

func (s *SuTaskAccountApi) Update(ctx context.Context, cmd *command.SuTaskAccountUpdateCommand) error {
	return s.accountService.Update(ctx, &cmd.Data)
}

func (s *SuTaskAccountApi) Delete(ctx context.Context, task *model2.SuTask) error {
	return nil
}

func (s *SuTaskAccountApi) BatchDelete(ctx context.Context, cmd *command.SuTaskAccountBatchDeleteCommand) error {
	return s.accountService.DeleteByIds(ctx, cmd.Data)
}

func (s *SuTaskAccountApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (store.FindPagingResult[*model2.SuTaskAccount], error) {
	res := s.accountService.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskAccountApi) FindById(ctx context.Context, qry *query.SuTaskAccountFindByIdQuery) (*model2.SuTaskAccount, error) {
	task, err := s.accountService.FindById(ctx, qry)
	return task, err
}
