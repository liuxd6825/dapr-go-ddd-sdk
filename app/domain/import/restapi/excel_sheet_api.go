package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type ExcelSheetApi struct {
	env      *env.Env
	rootPath string
	service  *service.ExcelSheetService
}

func NewExcelSheetApi(env *env.Env, rootPath string) *ExcelSheetApi {
	return &ExcelSheetApi{
		env:      env,
		rootPath: rootPath,
		service:  service.NewExcelSheetService(),
	}
}

func (s *ExcelSheetApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/import", "ExcelSheetApi", s)
	ctl.Post("sheet", "Create")
	ctl.Put("sheet", "Update")
	ctl.GetOne("sheet/{id}", "FindById")
	ctl.GetPaging("sheet", "FindPaging")
	ctl.GetData("sheet:file-id", "FindByFileId")
	ctl.GetData("sheet:doc-file-id", "FindByDocFileId")
	return ctl
}

func (s *ExcelSheetApi) Create(ctx context.Context, cmd *command.ExcelSheetCreateCommand) error {
	return s.service.Create(ctx, cmd.Data)
}

func (s *ExcelSheetApi) Update(ctx context.Context, cmd *command.ExcelSheetUpdateCommand) error {
	return s.service.Update(ctx, cmd.Data)
}

func (s *ExcelSheetApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.ExcelSheet], error) {
	return s.service.FindPaging(ctx, qry)
}

func (s *ExcelSheetApi) FindById(ctx context.Context, qry *query.ExcelSheetFindByIdQuery) (*model.ExcelSheet, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *ExcelSheetApi) FindByDocFileId(ctx context.Context, qry *query.ExcelSheetFindByDocFieldIdQuery) ([]*model.ExcelSheet, error) {
	return s.service.FindByDocFileId(ctx, qry.DocFileId)
}

func (s *ExcelSheetApi) FindByFileId(ctx context.Context, qry *query.ExcelSheetFindByFieldIdQuery) ([]*model.ExcelSheet, error) {
	return s.service.FindByFileId(ctx, qry.FileId)
}
