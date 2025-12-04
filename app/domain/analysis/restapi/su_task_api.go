package restapi

import (
	"context"

	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"

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

type SuTaskApi struct {
	taskService *service.SuTaskService
	codeService *service2.CodeService
	env         *env.Env
	rootPath    string
}

func NewSuTaskApi(env *env.Env, rootPath string) *SuTaskApi {
	taskService := service.NewSuTaskService()
	codeService := service2.NewCodeService()
	return &SuTaskApi{
		env:         env,
		rootPath:    rootPath,
		taskService: taskService,
		codeService: codeService,
	}
}

func (s *SuTaskApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "analysis.SuTaskApi", s)

	controller.Post("/analysis/su-task", "Create", restapi.WithTranDbKey(config.DBKey))
	controller.Put("/analysis/su-task", "Update", restapi.WithTranDbKey(config.DBKey))
	controller.Put("/analysis/su-task:analysis", "Analysis", restapi.WithTranDbKey(config.DBKey))

	controller.GetOne("/analysis/su-task/{id}", "FindById")
	controller.GetPaging("/analysis/su-task", "FindPaging")
	controller.GetOne("/analysis/su-task:bill/{id}", "FindBill")
	controller.GetData("/analysis/su-task:by-case", "FindByCaseId")

	controller.View("/analysis/su-task/bill.html", "GetBillView")

	return controller
}

func (s *SuTaskApi) Create(ctx context.Context, cmd *command.SuTaskCreateCommand) (*model2.SuTask, error) {
	return s.taskService.Create(ctx, cmd)
}

func (s *SuTaskApi) Update(ctx context.Context, cmd *command.SuTaskUpdateCommand) (*model2.SuTask, error) {
	return s.taskService.Update(ctx, cmd)
}

func (s *SuTaskApi) Analysis(ctx context.Context, cmd *command.SuTaskAnalysisCommand) error {
	return s.taskService.Analysis(ctx, cmd.Data.Id)
}

func (s *SuTaskApi) Delete(ctx context.Context, task *model2.SuTask) error {
	return nil
}

func (s *SuTaskApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (store.FindPagingResult[*model2.SuTask], error) {
	res := s.taskService.QueryPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskApi) FindById(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model2.SuTask, error) {
	task, err := s.taskService.QueryById(ctx, qry)
	return task, err
}

func (s *SuTaskApi) FindByCaseId(ctx context.Context, qry *query.SuTaskFindByCaseIdQuery) ([]*model2.SuTask, error) {
	tasks, err := s.taskService.FindByCaseId(ctx, qry.CaseId)
	return tasks, err
}

func (s *SuTaskApi) FindBill(ctx context.Context, qry *query.SuTaskFindByIdQuery) (*model2.SuTaskBillView, error) {
	task, err := s.taskService.QueryBillById(ctx, qry)
	return task, err
}

func (s *SuTaskApi) GetBillView(ctx context.Context, ictx iris.Context, qry *query.SuTaskFindByIdQuery) error {
	var task *model2.SuTask
	var err error
	if qry.Id != "" {
		task, err = s.FindById(ctx, qry)
		if err != nil {
			return err
		} else if task == nil {
			return iris.ErrNotFound
		}
	} else {
		task = model2.NewSuTask()
		task.Code = "newCode"
		task.Name = "新建任务"
	}

	isEditRule := false
	switch task.Status {
	case model2.SuTaskStatus_New:
		isEditRule = true
	default:
		isEditRule = false
	}

	vData := make(map[string]any)
	vData["task"] = task
	vData["isEditRule"] = isEditRule

	return restapi.View(ctx, ictx, "", vData)
}
