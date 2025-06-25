package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"io"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"

	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type DocumentAPI struct {
	env             *env.Env
	documentService *service.DocumentService
	fileService     *service.FileService
	fsService       *service.FsService
}

func NewDocumentAPI(env *env.Env, rootPath string) *DocumentAPI {
	documentService := service.NewDocumentService()
	fileService := service.NewFileService()
	fsService := service.NewFsService()
	return &DocumentAPI{
		env:             env,
		documentService: documentService,
		fileService:     fileService,
		fsService:       fsService,
	}
}

func (s *DocumentAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/doc/document:upload-chunk", "UploadChunk")
	b.Handle(iris.MethodGet, "/doc/document:download", "Download")
	b.Handle(iris.MethodPost, "/doc/document", "Create")
	b.Handle(iris.MethodPut, "/doc/document", "Update")
	b.Handle(iris.MethodPut, "/doc/document:rename", "Rename")
	b.Handle(iris.MethodDelete, "/doc/document", "Delete")
	b.Handle(iris.MethodGet, "/doc/document", "FindPaging")
}

func (s *DocumentAPI) UploadChunk(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		chunk, _, err := ictx.FormFile("chunk")
		chunkIndex := ictx.FormValue("chunkIndex")
		chunkSize := ictx.FormValue("chunkSize")
		objectName := ictx.FormValue("objectName")
		if err != nil {
			return err
		}

		has := s.fsService.Exists(objectName)

		if chunkIndex == "0" && !has {
			s.fsService.Create(objectName)
		}

		data, err := io.ReadAll(chunk)
		if err != nil {
			return err
		}

		err = s.fsService.WriteAt(objectName, data, chunkIndex, chunkSize)
		if err != nil {
			return err
		}

		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Download(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		fileId := ictx.URLParam("file-id")
		if fileId == "" {
			return errors.New("Url参数file-id不能为空")
		}

		file := s.fileService.FindById(ctx, fileId)
		if file == nil {
			return errors.New("没有找到文件记录,fileId=" + fileId)
		}

		has := s.fsService.Exists(file.ObjectName)
		if !has {
			return errors.New("没有找到文件,objectName=" + file.ObjectName)
		}

		err := s.fsService.Download(ictx, file.ObjectName, file.Name)
		if err != nil {
			return err
		}

		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.fileService.Create(ctx, s.fileService.GetFile(&cmd.Data))
		s.documentService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Rename(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"title"})
		s.documentService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.documentService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Delete(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentDeleteCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.documentService.DeleteById(ctx, cmd.Data.Id)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) FindPaging(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		folderId := ictx.URLParam("folder-id")
		//user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "folder_id=='" + folderId + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.documentService.FindPaging(ctx, qry)
		return web.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
