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

type ExcelFileApi struct {
	env         *env.Env
	rootPath    string
	fileService *service.ExcelFileService
}

func NewExcelFileApi(env *env.Env, rootPath string) *ExcelFileApi {
	fileService := service.NewExcelFileService()
	return &ExcelFileApi{
		env:         env,
		fileService: fileService,
		rootPath:    rootPath,
	}
}

func (s *ExcelFileApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/import", "ExcelFileApi", s)
	ctl.Post("file", "Create")
	ctl.Put("file", "Update")
	ctl.Delete("file", "Delete")
	ctl.GetPaging("file", "FindPaging")
	ctl.GetOne("file/{id}", "FindById")
	ctl.GetData("file:doc-file-id", "FindByDocFileId")
	return ctl
}

func (s *ExcelFileApi) Create(ctx context.Context, cmd *command.ExcelFileCreateCommand) (*model.ExcelFile, error) {
	return s.fileService.Create(ctx, cmd)
}

func (s *ExcelFileApi) Update(ctx context.Context, cmd *command.ExcelFileUpdateCommand) error {
	return s.fileService.Update(ctx, cmd)
}

func (s *ExcelFileApi) Delete(ctx context.Context, cmd *command.ExcelFileDeleteByIdCommand) error {
	return s.fileService.DeleteById(ctx, cmd)
}

func (s *ExcelFileApi) FindById(ctx context.Context, qry *query.ExcelFileFindByIdQuery) (*model.ExcelFile, error) {
	return s.fileService.FindById(ctx, qry.FileId)
}

func (s *ExcelFileApi) FindByDocFileId(ctx context.Context, qry *query.ExcelFileFindByDocFileIdQuery) (*model.ExcelFile, error) {
	return s.fileService.FindByDocFileId(ctx, qry.DocFileId)
}

func (s *ExcelFileApi) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (any, error) {
	res := s.fileService.FindPaging(ctx, qry)
	return res, res.GetError()
}
