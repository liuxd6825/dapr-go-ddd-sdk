package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordAPI struct {
	env           *env.Env
	rootPath      string
	recordService *service.RecordService
}

func NewRecordAPI(env *env.Env, rootPath string) *RecordAPI {
	return &RecordAPI{
		env:           env,
		rootPath:      rootPath,
		recordService: service.NewRecordService(),
	}
}

func (s *RecordAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath+"/import", s)
	ctl.Post("record", "Create")
	ctl.Post("record:generate-record", "GenerateRecord")
	ctl.Post("record:import2master", "Import2Master")
	ctl.Post("record:preview", "Preview")
	ctl.Delete("record", "Delete", restapi.WithParamsInBody(true))
	ctl.Post("record:revoke", "DeleteByTaskId")
	ctl.Put("record", "Update")
	ctl.Put("record:update-field", "UpdateField")
	ctl.Put("record:update-filter", "UpdateByFilter")
	ctl.GetOne("record/{id}", "FindById")
	ctl.GetPaging("record", "FindPaging")
	ctl.GetPaging("record:by-task", "FindPagingByTaskId")
	ctl.GetData("record:record-count", "CountRecordIeByTaskId")

	return nil
}

func (s *RecordAPI) Create(ctx context.Context, cmd *command.RecordCreateCommand) error {
	return s.recordService.Create(ctx, cmd)
}

func (s *RecordAPI) Update(ctx context.Context, cmd *command.RecordUpdateCommand) (any, error) {
	return s.recordService.Update(ctx, cmd)
}

func (s *RecordAPI) UpdateField(ctx context.Context, cmd *command.RecordUpdateFieldCommand) (*model.RecordIe, error) {
	res, _, err := s.recordService.UpdateField(ctx, cmd)
	return res, err
}

func (s *RecordAPI) UpdateByFilter(ctx context.Context, cmd *command.RecordUpdateFilterCommand) error {
	return s.recordService.UpdateByFilter(ctx, cmd)
}

func (s *RecordAPI) CountRecordIeByTaskId(ctx context.Context, qry *query.RecordIeFindByTaskIdQuery) (int64, error) {
	return s.recordService.CountRecordIeByTaskId(ctx, qry.TaskId)
}

func (s *RecordAPI) Delete(ctx context.Context, cmd *command.RecordDeleteCommand) error {
	return s.recordService.Delete(ctx, cmd)
}

func (s *RecordAPI) DeleteByTaskId(ctx context.Context, cmd *command.RecordDeleteCommand) error {
	return s.recordService.DeleteByTaskId(ctx, cmd.Data.TaskId)
}

func (s *RecordAPI) Preview(ctx context.Context, cmd *command.RecordPreviewCommand) ([]*model.RecordIe, error) {
	task := &model.Task{}
	task.Id = cmd.Data.TaskId
	task.CaseId = cmd.Data.CaseId
	task.DocId = cmd.Data.DocId
	task.FileId = cmd.Data.FileId
	task.FileName = cmd.Data.FileName
	task.SheetName = cmd.Data.SheetName

	res, _, err := s.recordService.Preview(ctx, task, cmd.Data.Template)
	return res, err
}

func (s *RecordAPI) GenerateRecord(ctx context.Context, cmd *command.RecordCreate4ExcelCommand) (*service.Create4ExcelResult, error) {
	return s.recordService.Create4Excel(ctx, cmd, func(batch readexcel.Batching) error {
		return nil
	})
}

func (s *RecordAPI) Import2Master(ctx context.Context, cmd *command.RecordImport2MasterCommand) error {
	return s.recordService.Import2Master(ctx, cmd)
}

func (s *RecordAPI) FindById(ctx context.Context, qry *query.RecordFindByIdQueryRequest) (*model.RecordIe, error) {
	return s.recordService.FindById(ctx, qry.Id)
}

func (s *RecordAPI) FindPagingByTaskId(ctx context.Context, qry *query.RecordIeFindPagingByTaskIdQuery) (store.FindPagingResult[*model.RecordIe], error) {
	res := s.recordService.FindPagingByTaskId(ctx, qry)
	return res, res.GetError()
}

func (s *RecordAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.RecordIe], error) {
	res := s.recordService.FindPaging(ctx, qry)
	return res, res.GetError()
}
