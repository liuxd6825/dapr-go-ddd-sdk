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
	service3 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CaseAPI struct {
	env           *env.Env
	caseService   *service.CaseService
	folderService *service2.FolderService
	fsService     *service2.FsService
	codeService   *service3.CodeService
	rootPath      string
}

func NewCaseAPI(env *env.Env, rootPath string) *CaseAPI {
	return &CaseAPI{
		env:           env,
		caseService:   service.NewCaseService(),
		folderService: service2.NewFolderService(),
		fsService:     service2.NewFsService(),
		codeService:   service3.NewCodeService(),
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

		code, err := s.codeService.New(ctx, "XM")
		if err != nil {
			return nil
		}

		cmd.Data.Code = code

		err = s.caseService.Create(ctx, cmd)
		if err != nil {
			return err
		}

		_, err = s.createRootFolder(ctx, cmd.Data.TenantId, cmd.Data.Id)
		if err != nil {
			return err
		}
		_, err = s.createMasterFolder(ctx, cmd.Data.TenantId, cmd.Data.Id)
		if err != nil {
			return err
		}

		hmFolder, err := s.createMasterTypeFolder(ctx, cmd.Data.TenantId, cmd.Data.Id, "人员")
		if err != nil {
			return err
		}
		cpFolder, err := s.createMasterTypeFolder(ctx, cmd.Data.TenantId, cmd.Data.Id, "公司")
		if err != nil {
			return err
		}
		htFolder, err := s.createMasterTypeFolder(ctx, cmd.Data.TenantId, cmd.Data.Id, "合同")
		if err != nil {
			return err
		}
		pdFolder, err := s.createMasterTypeFolder(ctx, cmd.Data.TenantId, cmd.Data.Id, "产品")
		if err != nil {
			return err
		}

		s.fsService.MkdirAll(hmFolder.FolderPath)
		s.fsService.MkdirAll(cpFolder.FolderPath)
		s.fsService.MkdirAll(htFolder.FolderPath)
		s.fsService.MkdirAll(pdFolder.FolderPath)
		return nil
	})
	return err
}

func (s *CaseAPI) createRootFolder(ctx context.Context, tenantId, caseId string) (*model2.Folder, error) {
	folder, _ := model2.NewFolder()
	folder.Id = tenantId + "_case_" + caseId
	folder.BusId = "case"
	folder.EntityId = caseId
	folder.RootId = tenantId + "_case_" + caseId
	folder.RootPath = "/" + tenantId + "/case/" + caseId
	folder.FolderPath = "/" + tenantId + "/case/" + caseId
	folder.CaseId = caseId
	folder.Name = "文件库"
	folder.DisabledFrontEdit = true

	err := s.folderService.CreateData(ctx, folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *CaseAPI) createMasterFolder(ctx context.Context, tenantId, caseId string) (*model2.Folder, error) {
	folder, _ := model2.NewFolder()
	folder.Id = tenantId + "_case_" + caseId + "_master"
	folder.BusId = "case"
	folder.EntityId = caseId
	folder.RootId = tenantId + "_case_" + caseId
	folder.RootPath = "/" + tenantId + "/case/" + caseId
	folder.FolderPath = "/" + tenantId + "/case/" + caseId + "/主数据附件"
	folder.CaseId = caseId
	folder.ParentId = tenantId + "_case_" + caseId
	folder.Name = "主数据附件"
	folder.DisabledFrontEdit = true

	err := s.folderService.CreateData(ctx, folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *CaseAPI) createMasterTypeFolder(ctx context.Context, tenantId, caseId, masterType string) (*model2.Folder, error) {
	var typeCode string
	switch masterType {
	case "人员":
		typeCode = "human"
		break
	case "公司":
		typeCode = "company"
		break
	case "合同":
		typeCode = "contract"
		break
	case "产品":
		typeCode = "product"
		break
	}
	folder, _ := model2.NewFolder()
	folder.Id = tenantId + "_case_" + caseId + "_master_" + typeCode
	folder.BusId = "case"
	folder.EntityId = caseId
	folder.RootId = tenantId + "_case_" + caseId
	folder.RootPath = "/" + tenantId + "/case/" + caseId
	folder.FolderPath = "/" + tenantId + "/case/" + caseId + "/主数据附件/" + masterType
	folder.CaseId = caseId
	folder.ParentId = tenantId + "_case_" + caseId + "_master"
	folder.Name = masterType
	folder.DisabledFrontEdit = true

	err := s.folderService.CreateData(ctx, folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
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
