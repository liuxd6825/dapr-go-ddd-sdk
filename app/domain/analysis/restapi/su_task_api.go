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

type SuTaskApi struct {
	taskService *service.SuTaskService
	env         *env.Env
	rootPath    string
}

func NewSuTaskApi(env *env.Env, rootPath string) *SuTaskApi {
	taskService := service.NewSuTaskService()
	return &SuTaskApi{
		env:         env,
		rootPath:    rootPath,
		taskService: taskService,
	}
}

func (s *SuTaskApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "analysis.SuTaskApi", s)
	controller.GetOne("/analysis/su-task/{id}", "FindById")
	controller.GetPaging("/analysis/su-task", "FindPaging")
	controller.Post("/analysis/su-task", "Create")
	controller.Put("/analysis/su-task", "Update")
	controller.Put("/analysis/su-task:start", "Start")
	controller.View("/analysis/su-task/bill.html", "GetBillView")
	controller.GetOne("/analysis/su-task:bill/{id}", "FindBill")
	return controller
}

func (s *SuTaskApi) Create(ctx context.Context, cmd *command.SuTaskCreateCommand) error {
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

func (s *SuTaskApi) Update(ctx context.Context, cmd *command.SuTaskCreateCommand) error {
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

func (s *SuTaskApi) Start(ctx context.Context, cmd *command.SuTaskStartCommand) error {
	return s.taskService.Start(ctx, cmd.Data.Id)
}

func (s *SuTaskApi) Delete(ctx context.Context, task *model2.SuTask) error {
	return nil
}

func (s *SuTaskApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (store.FindPagingResult[*model2.SuTask], error) {
	res := s.taskService.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskApi) FindById(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model2.SuTask, error) {
	task, err := s.taskService.FindById(ctx, qry)
	return task, err
}

func (s *SuTaskApi) FindBill(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model2.SuTaskBillView, error) {
	task, err := s.taskService.FindBillById(ctx, qry)
	return task, err
}

func (s *SuTaskApi) GetBillView(ctx context.Context, ictx iris.Context, qry *query.SuTaskFindByIdQuery) error {
	var task *model2.SuTask
	var err error
	if qry.Id != "" {
		task, err = s.FindById(ctx, qry)
		if err != nil {
			return err
		}
	} else {
		task = model2.NewSuTask()
		task.Code = "newCode"
		task.Name = "新建任务"
	}

	isEditRule := false
	switch task.Status {
	case model2.SuTaskStatus_New, "":
		isEditRule = true
	default:
		isEditRule = false
	}

	vData := make(map[string]any)
	vData["task"] = task
	vData["isEditRule"] = isEditRule

	return restapi.View(ctx, ictx, "", vData)
}
