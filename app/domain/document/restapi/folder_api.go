package restapi

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
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
}

func NewFolderAPI(env *env.Env, rootPath string) *FolderAPI {
	folderService := service.NewFolderService()
	docService := service.NewDocumentService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	return &FolderAPI{
		env:           env,
		folderService: folderService,
		docService:    docService,
		fileService:   fileService,
		fsService:     fsService,
	}
}

func (s *FolderAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/doc/folder:root", "CreateRoot")
	b.Handle(iris.MethodPost, "/doc/folder", "Create")
	b.Handle(iris.MethodPut, "/doc/folder", "Update")
	b.Handle(iris.MethodPut, "/doc/folder:rename", "Rename")
	b.Handle(iris.MethodPut, "/doc/folder:color", "SetColor")
	b.Handle(iris.MethodPut, "/doc/folder:move", "Move")
	b.Handle(iris.MethodDelete, "/doc/folder", "Delete")
	b.Handle(iris.MethodGet, "/doc/folder", "FindPaging")
}

func (s *FolderAPI) CreateRoot(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FolderCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

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
			cmd.Data.Name = "根目录"
			s.folderService.Create(ctx, &cmd.Data)
			s.fsService.MkdirAll(cmd.Data.FolderPath)
			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FolderCreateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

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

			s.folderService.Create(ctx, &cmd.Data)
			s.fsService.MkdirAll(cmd.Data.FolderPath)
			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) Rename(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FolderRenameCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			res, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
			if err != nil {
				return err
			}
			folders := []*model.Folder{}
			for _, f := range res {
				if strings.HasPrefix(f.FolderPath+"/", cmd.Data.OldName+"/") {
					f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.OldName+"/", cmd.Data.FolderPath+"/", -1)
					f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
					folders = append(folders, f)
				}
			}
			if len(folders) > 0 {
				opts := idao.NewCallOptions()
				opts.SetUpdateFields([]string{"folder_path"})
				s.folderService.UpdateMany(ctx, folders, opts)
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"name", "folder_path"})

			folder := model.Folder{}
			folder.Id = cmd.Data.Id
			folder.Name = cmd.Data.Name
			folder.FolderPath = cmd.Data.FolderPath

			s.folderService.Update(ctx, &folder, opts)

			err = s.fsService.Rename(cmd.Data.OldName, cmd.Data.FolderPath)
			if err != nil {
				return err
			}
			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) SetColor(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"color"})
		s.folderService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) Move(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FolderMoveCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			res, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
			if err != nil {
				return err
			}
			folders := []*model.Folder{}
			for _, f := range res {
				if strings.HasPrefix(f.FolderPath+"/", cmd.Data.SourcePath+"/") {
					f.FolderPath = strings.Replace(f.FolderPath+"/", cmd.Data.SourcePath+"/", cmd.Data.TargetPath+"/"+cmd.Data.Name+"/", -1)
					f.FolderPath = strings.TrimSuffix(f.FolderPath, "/")
					folders = append(folders, f)
				}
			}
			if len(folders) > 0 {
				opts := idao.NewCallOptions()
				opts.SetUpdateFields([]string{"folder_path"})
				s.folderService.UpdateMany(ctx, folders, opts)
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"parent_id"})

			folder := model.Folder{}
			folder.Id = cmd.Data.Id
			folder.ParentId = cmd.Data.ParentId

			s.folderService.Update(ctx, &folder, opts)

			err = s.fsService.MoveDir(cmd.Data.SourcePath, cmd.Data.TargetPath+"/"+cmd.Data.Name)
			if err != nil {
				return err
			}
			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.folderService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.folderService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FolderDeleteCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			//has := s.folderService.HasChildren(ctx, cmd.Data.Id)
			//if has {
			//	return errors.New("存在子目录")
			//}
			//has = s.docService.HasDocumentByFolder(ctx, cmd.Data.Id)
			//if has {
			//	return errors.New("存在文档")
			//}

			res, err := s.folderService.FindByRSQL(ctx, fmt.Sprintf("tenant_id==\"%s\" and bus_id==\"%s\" and entity_id==\"%s\"", cmd.Data.TenantId, cmd.Data.BusId, cmd.Data.EntityId))
			if err != nil {
				return err
			}
			for _, f := range res {
				if strings.HasPrefix(f.FolderPath+"/", cmd.Data.FolderPath+"/") {
					s.fileService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
					s.docService.DeleteByRSQL(ctx, fmt.Sprintf("folder_id=='%s'", f.Id))
					s.folderService.DeleteById(ctx, f.Id)
				}
			}
			s.fsService.RemoveAll(cmd.Data.FolderPath)
			return nil
		})

		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FolderAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		parentId := ictx.URLParam("folder-id")
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "parent_id=='" + parentId + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.folderService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
