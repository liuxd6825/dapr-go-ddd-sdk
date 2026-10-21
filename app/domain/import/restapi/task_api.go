package restapi

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TaskAPI struct {
	dbKey         string
	env           *envs.Env
	taskService   *service.TaskService
	recordService *service.RecordService
	rootPath      string
}

func NewTaskAPI(env *envs.Env, rootPath string, livyServer *livyserver.LivyServer, dbKey string) *TaskAPI {
	return &TaskAPI{
		dbKey:         dbKey,
		env:           env,
		rootPath:      rootPath,
		taskService:   service.NewTaskService(dbKey),
		recordService: service.NewRecordService(livyServer, dbKey),
	}
}

func (s *TaskAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.taskService = service.NewTaskService(s.dbKey)
	ctl := restapi.NewController(app, s.rootPath+"/import", "TaskAPI", s)
	ctl.Post("/task", "Create")
	ctl.Put("/task", "Update")
	ctl.Delete("/task", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/task/{id}", "FindById")
	ctl.GetOne("/task:by-file-id", "FindByFileId")
	ctl.GetPaging("/task", "FindPaging")
	return ctl
}

func (s *TaskAPI) Create(ctx context.Context, cmd *command.TaskCreateCommand) error {
	return tx.StartTx(ctx, tx.NewTxCfg(s.dbKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		return s.taskService.Create(ctx, cmd)
	})
}

func (s *TaskAPI) Update(ctx context.Context, cmd *command.TaskUpdateCommand) error {
	return s.taskService.Update(ctx, cmd)
}

func (s *TaskAPI) Delete(ctx context.Context, cmd *command.TaskDeleteCommand) error {
	return s.taskService.Delete(ctx, cmd)
}

func (s *TaskAPI) FindById(ctx context.Context, qry *query.TaskFindByIdQuery) (*model.Task, error) {
	return s.taskService.FindById(ctx, qry.Id)
}

func (s *TaskAPI) FindPaging(ctx context.Context, qry *idao.FindPagingByCaseIdQueryRequest) (idao.FindPagingResult[*model.Task], error) {
	return s.taskService.FindPagingByCaseId(ctx, qry)
}

func (s *TaskAPI) FindByFileId(ctx context.Context, qry *query.TaskFindByFileIdQuery) (*model.Task, error) {
	return s.taskService.FindByFileId(ctx, qry.CaseId, qry.FileId)
}
