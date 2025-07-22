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
	ctl.Post("record:import2master", "Import2Master")
	ctl.Delete("record", "Delete")
	ctl.Put("record", "Update")
	ctl.GetOne("record/{id}", "FindById")
	ctl.GetPaging("record", "FindPaging")
	return nil
}

func (s *RecordAPI) Create(ctx context.Context, cmd *command.RecordCreateCommand) error {
	return s.recordService.Create(ctx, cmd)
}

func (s *RecordAPI) Update(ctx context.Context, cmd *command.RecordUpdateCommand) (any, error) {
	return s.recordService.Update(ctx, cmd)
}

func (s *RecordAPI) Delete(ctx context.Context, cmd *command.RecordDeleteCommand) error {
	return s.recordService.Delete(ctx, cmd)
}

func (s *RecordAPI) Import2Master(ctx context.Context, cmd *command.RecordImport2MasterCommand) error {
	return s.recordService.Import2Master(ctx, cmd)
}

func (s *RecordAPI) FindById(ctx context.Context, qry *query.RecordFindByIdQueryRequest) (*model.RecordIe, error) {
	return s.recordService.FindById(ctx, qry.Id)
}

func (s *RecordAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.RecordIe], error) {
	res := s.recordService.FindPaging(ctx, qry)
	return res, res.GetError()
}
