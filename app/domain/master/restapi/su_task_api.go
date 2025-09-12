package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
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
	controller := restapi.NewController(app, s.rootPath, "master.SuTaskApi", s)
	controller.GetOne("/su-task/{id}", "FindById")
	controller.GetPaging("/su-task", "FindPaging")
	controller.Post("/su-task", "Create")
	controller.Put("/su-task", "Update")
	controller.Post("/su-task:analyse", "Analyse")
	controller.View("/su/task/bill.html", "BillView")
	return controller
}

func (s *SuTaskApi) Create(ctx context.Context, task *model.SuTask) error {
	return s.taskService.Create(ctx, task)
}

func (s *SuTaskApi) Update(ctx context.Context, task *model.SuTask) error {
	return s.taskService.Update(ctx, task)
}

func (s *SuTaskApi) Analyse(ctx context.Context, task *model.SuTask) error {
	return nil
}

func (s *SuTaskApi) Delete(ctx context.Context, task *model.SuTask) error {
	return nil
}

func (s *SuTaskApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (store.FindPagingResult[*model.SuTask], error) {
	res := s.taskService.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskApi) FindById(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model.SuTask, error) {
	return s.taskService.FindById(ctx, qry)
}

func (s *SuTaskApi) BillView(ctx context.Context, ictx iris.Context, qry *query.SuTaskFindByIdQuery) error {
	var task *model.SuTask
	var err error
	if qry.Id != "" {
		task, err = s.FindById(ctx, qry)
		if err != nil {
			return err
		}
	} else {
		task = model.NewSuTask()
		task.Code = "newCode"
		task.Name = "新建任务"
	}

	isEditRule := false
	switch task.Status {
	case model.SuTaskStatus_New, "":
		isEditRule = true
	default:
		isEditRule = false
	}

	vData := make(map[string]any)
	vData["task"] = task
	vData["isEditRule"] = isEditRule

	return restapi.View(ctx, ictx, "", vData)
}
