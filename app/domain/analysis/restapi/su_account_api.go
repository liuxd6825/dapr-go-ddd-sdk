package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/command"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type SuTaskAccountApi struct {
	taskService *service.SuTaskService
	env         *env.Env
	rootPath    string
}

func NewSuTaskAccountApi(env *env.Env, rootPath string) *SuTaskAccountApi {
	taskService := service.NewSuTaskService()
	return &SuTaskAccountApi{
		env:         env,
		rootPath:    rootPath,
		taskService: taskService,
	}
}

func (s *SuTaskAccountApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.SuAccountApi", s)
	controller.Post("su-account", "Create")
	controller.Put("su-account", "Update")
	controller.GetOne("su-account/{id}", "FindById")
	controller.GetPaging("su-account", "FindPaging")
	return controller
}

func (s *SuTaskAccountApi) Create(ctx context.Context, cmd *command.SuTaskCreateCommand) error {
	task := &model2.SuTask{
		BaseModel:   cmd.Data.BaseModel,
		Code:        cmd.Data.Code,
		Name:        cmd.Data.Name,
		StartTime:   cmd.Data.StartTime,
		EndTime:     cmd.Data.EndTime,
		Status:      model2.SuTaskStatus_New,
		Rules:       cmd.Data.Rules,
		SuCount:     0,
		SuHighCount: 0,
	}
	return s.taskService.Create(ctx, task)
}

func (s *SuTaskAccountApi) Update(ctx context.Context, cmd *command.SuTaskCreateCommand) error {
	task := &model2.SuTask{
		BaseModel:   cmd.Data.BaseModel,
		Code:        cmd.Data.Code,
		Name:        cmd.Data.Name,
		StartTime:   cmd.Data.StartTime,
		EndTime:     cmd.Data.EndTime,
		Status:      model2.SuTaskStatus_New,
		Rules:       cmd.Data.Rules,
		SuCount:     0,
		SuHighCount: 0,
	}
	return s.taskService.Update(ctx, task)
}

func (s *SuTaskAccountApi) Delete(ctx context.Context, task *model2.SuTask) error {
	return nil
}

func (s *SuTaskAccountApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (store.FindPagingResult[*model2.SuTask], error) {
	res := s.taskService.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskAccountApi) FindById(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model2.SuTask, error) {
	task, err := s.taskService.FindById(ctx, qry)
	return task, err
}
