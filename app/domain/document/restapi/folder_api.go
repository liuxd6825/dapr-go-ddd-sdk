package restapi

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"strings"
)

type FolderAPI struct {
	env           *env.Env
	folderService *service.FolderService
	docService    *service.DocumentService
	fileService   *service.FileService
	fsService     *service.FsService
	rootPath      string
}

func NewFolderAPI(env *env.Env, rootPath string) *FolderAPI {
	folderService := service.NewFolderService()
	docService := service.NewDocumentService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	return &FolderAPI{
		rootPath:      rootPath,
		env:           env,
		folderService: folderService,
		docService:    docService,
		fileService:   fileService,
		fsService:     fsService,
	}
}

func (s *FolderAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.FolderAPI", s)
	ctl.Post("/folder:root", "CreateRoot")
	ctl.Post("/folder", "Create")
	ctl.Put("/folder", "Update")
	ctl.Put("/folder", "Update")
	ctl.Put("/folder:rename", "Rename")
	ctl.Put("/folder:color", "SetColor")
	ctl.Put("/folder:move", "Move")
	ctl.Delete("/folder", "Delete", restapi.WithParamsInBody(true))
	ctl.GetPaging("/folder", "FindPaging")
	ctl.GetOne("/folder/:id", "FindById")
	ctl.GetData("/folder:tree", "FindTree")
	return ctl
}

func (s *FolderAPI) CreateRoot(ctx context.Context, cmd *command.FolderCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {

		vErr := errors.NewVerifyError()
		if cmd.Data.TenantId == "" {
			vErr.AppendField("tenantId", "不能为空", "租户Id")
		}
		if cmd.Data.BusId == "" {
			vErr.AppendField("busId", "不能为空", "业务Id")
		}
		if cmd.Data.EntityId == "" {
			vErr.AppendField("entityId", "不能为空", "实体Id")
		}
		if vErr.HasError() {
			return vErr
		}

		cmd.Data.Id = cmd.Data.TenantId + "_" + cmd.Data.BusId + "_" + cmd.Data.EntityId
		cmd.Data.RootId = cmd.Data.TenantId + "_" + cmd.Data.BusId + "_" + cmd.Data.EntityId
		cmd.Data.RootPath = "/" + cmd.Data.TenantId + "/" + cmd.Data.BusId + "/" + cmd.Data.EntityId
		cmd.Data.FolderPath = "/" + cmd.Data.TenantId + "/" + cmd.Data.BusId + "/" + cmd.Data.EntityId
		cmd.Data.Name = "文件库"
		err1 := s.folderService.Create(ctx, cmd)
		if err1 != nil {
			return err1
		}

		s.fsService.MkdirAll(cmd.Data.FolderPath)
		return nil
	})
	return err
}

func (s *FolderAPI) Create(ctx context.Context, cmd *command.FolderCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {

		vErr := errors.NewVerifyError()
		if cmd.Data.TenantId == "" {
			vErr.AppendField("tenantId", "不能为空", "租户Id")
		}
		if cmd.Data.BusId == "" {
			vErr.AppendField("busId", "不能为空", "业务Id")
		}
		if cmd.Data.EntityId == "" {
			vErr.AppendField("entityId", "不能为空", "实体Id")
		}
		if vErr.HasError() {
			return vErr
		}

		err1 := s.folderService.Create(ctx, cmd)
		if err1 != nil {
			return err1
		}

		s.fsService.MkdirAll(cmd.Data.FolderPath)
		return nil
	})

	return err
}

func (s *FolderAPI) Rename(ctx context.Context, cmd *command.FolderRenameCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {

		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}
		folders := []*model.Folder{}
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.OldName+"/") {
				f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.OldName+"/", cmd.Data.FolderPath+"/", -1)
				f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
				folders = append(folders, f)
			}
		}

		if len(folders) > 0 {
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"folder_path"})
			err = s.folderService.UpdateMany(ctx, folders, opts)
			if err != nil {
				return err
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "folder_path"})

		folder := model.Folder{}
		folder.Id = cmd.Data.Id
		folder.Name = cmd.Data.Name
		folder.FolderPath = cmd.Data.FolderPath

		err = s.folderService.Update(ctx, &folder, opts)
		if err != nil {
			return err
		}

		err = s.fsService.Rename(cmd.Data.OldName, cmd.Data.FolderPath)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *FolderAPI) SetColor(ctx context.Context, cmd *command.FolderUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"color"})
	err := s.folderService.Update(ctx, &cmd.Data, opts)
	return err
}

func (s *FolderAPI) Move(ctx context.Context, cmd *command.FolderMoveCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {

		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}
		folders := []*model.Folder{}
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.SourcePath+"/") {
				f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.SourcePath+"/", cmd.Data.TargetPath+"/"+cmd.Data.Name+"/", -1)
				f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
				folders = append(folders, f)
			}
		}

		if len(folders) > 0 {
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"folder_path"})
			err = s.folderService.UpdateMany(ctx, folders, opts)
			if err != nil {
				return err
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"parent_id"})

		folder := model.Folder{}
		folder.Id = cmd.Data.Id
		folder.ParentId = cmd.Data.ParentId

		err = s.folderService.Update(ctx, &folder, opts)
		if err != nil {
			return err
		}

		err = s.fsService.MoveDir(cmd.Data.SourcePath, cmd.Data.TargetPath+"/"+cmd.Data.Name)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *FolderAPI) Update(ctx context.Context, cmd *command.FolderUpdateCommand) error {
	err := s.folderService.Update(ctx, &cmd.Data)
	return err
}

func (s *FolderAPI) Delete(ctx context.Context, cmd *command.FolderDeleteCommand) error {
	err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {

		arr, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
		if err != nil {
			return err
		}

		var res *idao.Result
		for _, f := range arr {
			if strings.HasPrefix(f.FolderPath+"/", cmd.Data.FolderPath+"/") {
				res = s.fileService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
				if res.Error != nil {
					return res.Error
				}

				res = s.docService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
				if res.Error != nil {
					return res.Error
				}

				res = s.folderService.DeleteById(ctx, f.Id)
				if res.Error != nil {
					return res.Error
				}
			}
		}
		s.fsService.RemoveAll(cmd.Data.FolderPath)
		return nil
	})

	return err
}

func (s *FolderAPI) FindPaging(ctx context.Context, query *query.FindFolderByFolderIdQuery) (idao.FindPagingResult[*model.Folder], error) {
	qry := store.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = "parent_id=='" + query.FolderId + "'"
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	return s.folderService.FindPaging(ctx, qry)
}

func (s *FolderAPI) FindTree(ctx context.Context, query *query.FindTreeByParamQuery) []model.FolderTree {
	qry := store.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = fmt.Sprintf("tenant_id=='%s' and bus_id=='%s' and entity_id=='%s'", query.TenantId, query.BusId, query.EntityId)
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	res, _ := s.folderService.FindPaging(ctx, qry)
	return s.folderService.TranslateTreeData(res.GetData())
}

func (s *FolderAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Folder, error) {
	return s.folderService.FindById(ctx, qry.Id)
}
