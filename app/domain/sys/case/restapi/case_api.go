package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CaseAPI struct {
	env           *env.Env
	caseService   *service.CaseService
	folderService *service2.FolderService
	fsService     *service2.FsService
	rootPath      string
}

func NewCaseAPI(env *env.Env, rootPath string) *CaseAPI {
	return &CaseAPI{
		env:           env,
		caseService:   service.NewCaseService(),
		folderService: service2.NewFolderService(),
		fsService:     service2.NewFsService(),
		rootPath:      rootPath,
	}
}

func (s *CaseAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.caseService = service.NewCaseService()
	ctl := restapi.NewController(app, s.rootPath+"/master", "sys.CaseApi", s)
	ctl.Post("/case", "Create")
	ctl.Put("/case", "Update")
	ctl.Delete("/case", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/case:deleteBatch", "DeleteBatch", restapi.WithParamsInBody(true))
	ctl.GetOne("/case/{id}", "FindById")
	ctl.GetPaging("/case", "FindPaging")
	return ctl
}

func (s *CaseAPI) Create(ctx context.Context, cmd *command.CaseCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.caseService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		err := s.caseService.Create(ctx, cmd)
		if err != nil {
			return err
		}

		folder, _ := model2.NewFolder()
		folder.Id = cmd.Data.TenantId + "_case_" + cmd.Data.Id
		folder.RootId = cmd.Data.TenantId + "_case_" + cmd.Data.Id
		folder.RootPath = "/" + cmd.Data.TenantId + "/case/" + cmd.Data.Id
		folder.FolderPath = "/" + cmd.Data.TenantId + "/case/" + cmd.Data.Id
		folder.CaseId = cmd.Data.Id
		folder.Name = "文件库"

		err = s.folderService.CreateData(ctx, folder)
		if err != nil {
			return err
		}

		s.fsService.MkdirAll(folder.FolderPath)
		return nil
	})
	return err
}

func (s *CaseAPI) Update(ctx context.Context, cmd *command.CaseUpdateCommand) error {
	return s.caseService.Update(ctx, cmd)
}

func (s *CaseAPI) Delete(ctx context.Context, cmd *command.CaseDeleteCommand) error {
	return s.caseService.Delete(ctx, cmd)
}

func (s *CaseAPI) DeleteBatch(ctx context.Context, cmd *command.CaseDeleteBatchCommand) error {
	return s.caseService.DeleteBatch(ctx, cmd)
}

func (s *CaseAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Case, error) {
	return s.caseService.FindById(ctx, qry)
}

func (s *CaseAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.Case], error) {
	return s.caseService.FindPaging(ctx, qry)
}
