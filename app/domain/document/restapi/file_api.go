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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type FileAPI struct {
	env         *env.Env
	fileService *service.FileService
	docService  *service.DocumentService
}

func NewFileAPI(env *env.Env, rootPath string) *FileAPI {
	fileService := service.NewFileService()
	docService := service.NewDocumentService()
	return &FileAPI{
		env:         env,
		fileService: fileService,
		docService:  docService,
	}
}

func (s *FileAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/doc/file", "Create")
	b.Handle(iris.MethodPut, "/doc/file:ver-info", "UpdateVerInfo")
	b.Handle(iris.MethodPut, "/doc/file:main", "UpdateIsMain")
	b.Handle(iris.MethodPut, "/doc/file", "Update")
	b.Handle(iris.MethodDelete, "/doc/file", "Delete")
	b.Handle(iris.MethodGet, "/doc/file", "FindPaging")
	b.Handle(iris.MethodGet, "/doc/file:document-id", "FindByDocumentId")
}

func (s *FileAPI) Create(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FileCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.fileService.Create(ctx, &cmd.Data)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) UpdateVerInfo(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FileUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"ver_info"})
		res := s.fileService.Update(ctx, &cmd.Data, opts)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) UpdateIsMain(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		err := tx.StartTx(ctx, []string{s.fileService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
			var cmd *command.FileUpdateCommand
			if err := ictx.ReadJSON(&cmd); err != nil {
				return err
			}

			files, err := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", cmd.Data.DocumentId))
			if err != nil {
				return err
			}

			var res *idao.Result

			if len(files) > 0 {
				for _, file := range files {
					file.IsMain = false
				}
				opts := idao.NewCallOptions()
				opts.SetUpdateFields([]string{"is_main"})
				res = s.fileService.UpdateMany(ctx, files, opts)
				if res.Error != nil {
					return res.Error
				}
			}

			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"is_main"})
			res = s.fileService.Update(ctx, &cmd.Data, opts)
			if res.Error != nil {
				return res.Error
			}

			opts.SetUpdateFields([]string{"object_name", "name", "download_url", "preview_url", "file_id", "size_title", "size", "ext_name", "thumbnail"})
			doc := model.Document{}
			doc.Id = cmd.Data.DocumentId
			doc.ObjectName = cmd.Data.ObjectName
			doc.Name = cmd.Data.Name
			doc.DownloadUrl = cmd.Data.DownloadUrl
			doc.PreviewUrl = cmd.Data.PreviewUrl
			doc.FileId = cmd.Data.Id
			doc.SizeTitle = cmd.Data.SizeTitle
			doc.Size = cmd.Data.Size
			doc.ExtName = cmd.Data.ExtName
			doc.Thumbnail = cmd.Data.Thumbnail
			res = s.docService.Update(ctx, &doc, opts)
			if res.Error != nil {
				return res.Error
			}

			return nil
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) Update(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FileUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.fileService.Update(ctx, &cmd.Data)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) Delete(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		var cmd *command.FileDeleteCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		res := s.fileService.DeleteById(ctx, cmd.Data.Id)
		return res.Error
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) FindPaging(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.URLParam("case-id")
		user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "case_id=='" + caseId + "' and creator_id=='" + user.GetId() + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.fileService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}

func (s *FileAPI) FindByDocumentId(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		docId := ictx.URLParam("document-id")
		files, err := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", docId))
		if err != nil {
			return err
		}
		return restapi.SetData(ictx, files)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
