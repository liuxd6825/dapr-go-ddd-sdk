package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"

	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type FolderAPI struct {
	env           *env.Env
	folderService *service.FolderService
	docService    *service.DocumentService
}

func NewFolderAPI(env *env.Env, rootPath string) *FolderAPI {
	folderService := service.NewFolderService()
	docService := service.NewDocumentService()
	return &FolderAPI{
		env:           env,
		folderService: folderService,
		docService:    docService,
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
	web.Try(ictx, func(ctx context.Context) error {
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
		cmd.Data.Name = "根目录"

		s.folderService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
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
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) Rename(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"name", "updatedTime", "updaterId", "updaterName"})
		s.folderService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) SetColor(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"color", "updatedTime", "updaterId", "updaterName"})
		s.folderService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) Move(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"parentId", "updatedTime", "updaterId", "updaterName"})
		s.folderService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.folderService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) Delete(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FolderDeleteCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}

		has := s.folderService.HasChildren(ctx, cmd.Data.Id)
		if has {
			return errors.New("存在子目录")
		}
		has = s.docService.HasDocumentByFolder(ctx, cmd.Data.Id)
		if has {
			return errors.New("存在文档")
		}

		s.folderService.DeleteById(ctx, cmd.Data.Id)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *FolderAPI) FindPaging(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		parentId := ictx.URLParam("folder-id")
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "parent_id=='" + parentId + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.folderService.FindPaging(ctx, qry)
		return web.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
