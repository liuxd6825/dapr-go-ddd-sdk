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

type ExcelRowApi struct {
	env        *env.Env
	rootPath   string
	rowService *service.ExcelRowService
}

func NewExcelRowApi(env *env.Env, rootPath string) *ExcelRowApi {
	return &ExcelRowApi{
		env:        env,
		rootPath:   rootPath,
		rowService: service.NewExcelRowService(),
	}
}

func (s *ExcelRowApi) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath+"/import", s)
	ctl.Post("row", "Create")
	ctl.Put("row", "Update")
	ctl.Delete("row", "Delete")
	ctl.GetPaging("row", "FindPaging")
	ctl.GetOne("row/{id}", "FindById")
	ctl.GetData("row:rows", "FindRows")
	return nil
}

func (s *ExcelRowApi) Create(ctx context.Context, cmd *command.ExcelRowCreateCommand) error {
	return s.rowService.Create(ctx, cmd)
}

func (s *ExcelRowApi) Update(ctx context.Context, cmd *command.ExcelRowUpdateCommand) error {
	return s.rowService.Update(ctx, cmd.Data, cmd.UpdateMask)
}

func (s *ExcelRowApi) Delete(ctx context.Context, cmd *command.ExcelRowDeleteCommand) error {
	return s.rowService.DeleteById(ctx, cmd.Data.Id)
}

func (s *ExcelRowApi) FindById(ctx context.Context, qry *query.ExcelRowFindByIdQuery) (*model.ExcelRow, error) {
	return s.rowService.FindById(ctx, qry.Id)
}

func (s *ExcelRowApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.ExcelRow], error) {
	res := s.rowService.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *ExcelRowApi) FindRows(ctx context.Context, qry *query.ExcelRowFindPreviewQuery) (*query.ExcelRowFindPreviewQueryResult, error) {
	return s.rowService.FindRows(ctx, qry)
}
