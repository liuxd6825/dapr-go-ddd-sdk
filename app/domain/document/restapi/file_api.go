package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type FileAPI struct {
	rootPath    string
	env         *env.Env
	fileService *service.FileService
	docService  *service.DocumentService
}

func NewFileAPI(env *env.Env, rootPath string) *FileAPI {
	fileService := service.NewFileService()
	docService := service.NewDocumentService()
	return &FileAPI{
		rootPath:    rootPath,
		env:         env,
		fileService: fileService,
		docService:  docService,
	}
}

func (s *FileAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.FileAPI", s)
	ctl.Post("/file", "Create")
	ctl.Put("/file:ver-info", "UpdateVerInfo")
	ctl.Put("/file:main", "UpdateIsMain")
	ctl.Put("/file", "Update")
	ctl.Delete("/file", "Delete", restapi.WithParamsInBody(true))
	ctl.GetData("/file", "FindPaging")
	ctl.GetData("/file:document-id", "FindByDocumentId")
	return ctl
}

func (s *FileAPI) Create(ctx context.Context, cmd *command.FileCreateCommand) error {
	return s.fileService.Create(ctx, &cmd.Data)
}

func (s *FileAPI) UpdateVerInfo(ctx context.Context, cmd *command.FileUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"ver_info"})
	return s.fileService.Update(ctx, &cmd.Data, opts)
}

func (s *FileAPI) UpdateIsMain(ctx context.Context, cmd *command.FileUpdateCommand) error {
	err := tx.StartTx(ctx, []string{s.fileService.GetConfig().DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		files, err := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", cmd.Data.DocumentId))
		if err != nil {
			return err
		}

		if len(files) > 0 {
			for _, file := range files {
				file.IsMain = false
			}
			opts := idao.NewCallOptions()
			opts.SetUpdateFields([]string{"is_main"})
			err = s.fileService.UpdateMany(ctx, files, opts)
			if err != nil {
				return err
			}
		}

		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"is_main"})
		err = s.fileService.Update(ctx, &cmd.Data, opts)
		if err != nil {
			return err
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
		err = s.docService.Update(ctx, &doc, opts)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *FileAPI) Update(ctx context.Context, cmd *command.FileUpdateCommand) error {
	return s.fileService.Update(ctx, &cmd.Data)
}

func (s *FileAPI) Delete(ctx context.Context, cmd *command.FileDeleteCommand) error {
	return s.fileService.DeleteById(ctx, cmd.Data.Id).GetError()
}

func (s *FileAPI) FindPaging(ctx context.Context, query *query.FindByCaseIdQuery) (idao.FindPagingResult[*model.File], error) {
	user, _ := appctx.GetAuthUser(ctx)
	qry := store2.NewFindPagingQueryRequest()
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = "case_id=='" + query.CaseId + "' and creator_id=='" + user.GetId() + "'"
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	return s.fileService.FindPaging(ctx, qry)
}

func (s *FileAPI) FindByDocumentId(ctx context.Context, query *query.FindByDocumentIdQuery) []*model.File {
	files, _ := s.fileService.FindByRSQL(ctx, fmt.Sprintf("document_id=='%s'", query.DocumentId))
	return files
}
