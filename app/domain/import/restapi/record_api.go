package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TaskAPI struct {
	env           *env.Env
	recordService *service.RecordService
	rootPath      string
}

func NewTaskAPI(env *env.Env, rootPath string) *TaskAPI {
	return &TaskAPI{
		env:           env,
		recordService: service.NewRecordService(),
		rootPath:      rootPath,
	}
}

func (s *TaskAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath, s)
	ctl.Post("/tasks", "Create")
	ctl.Put("/tasks", "Update")
	ctl.GetOne("/tasks/{id}", "FindById")
	ctl.GetPaging("/tasks", "FindPaging")
	return nil
}

func (s *TaskAPI) Create(ctx context.Context, cmd *command.RecordCreateCommand) error {
	return s.recordService.Create(ctx, cmd)
}

func (s *TaskAPI) Update(ctx context.Context, cmd *command.RecordUpdateCommand) (any, error) {
	return s.recordService.Update(ctx, cmd)
}

func (s *TaskAPI) Delete(ctx context.Context, cmd *command.RecordDeleteCommand) error {
	return s.recordService.Delete(ctx, cmd)
}

func (s *TaskAPI) FindById(ctx context.Context, qry *query.RecordFindByIdQueryRequest) (*model.RecordIe, error) {
	return s.recordService.FindById(ctx, qry.Id)
}

func (s *TaskAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.RecordIe], error) {
	res := s.recordService.FindPaging(ctx, qry)
	return res, res.GetError()
}
