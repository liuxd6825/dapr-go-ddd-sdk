package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentMetaAPI struct {
	rootPath            string
	env                 *env.Env
	documentMetaService *service.DocumentMetaService
	docService          *service.DocumentService
	folderService       *service.FolderService
}

func NewDocumentMetaAPI(env *env.Env, rootPath string) *DocumentMetaAPI {
	documentMetaService := service.NewDocumentMetaService()
	docService := service.NewDocumentService()
	folderService := service.NewFolderService()
	return &DocumentMetaAPI{
		rootPath:            rootPath,
		env:                 env,
		documentMetaService: documentMetaService,
		docService:          docService,
		folderService:       folderService,
	}
}

func (s *DocumentMetaAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.DocumentMetaAPI", s)
	ctl.Post("/document-meta", "Create")
	ctl.Post("/document-meta:save", "Save")
	ctl.Post("/document-meta:batch", "CreateMany")
	ctl.Put("/document-meta", "Update")
	return ctl
}

func (s *DocumentMetaAPI) Create(ctx context.Context, cmd *command.DocumentMetaCreateCommand) error {
	err := s.documentMetaService.Create(ctx, &cmd.Data)
	return err
}

func (s *DocumentMetaAPI) Save(ctx context.Context, cmd *command.DocumentMetaSubmitCommand) error {
	err := tx.StartTx(ctx, []string{config.DBKey}, func(ctx context.Context, options ...*store2.SessionOptions) error {
		err := s.documentMetaService.Save(ctx, cmd)
		if err != nil {
			return err
		}
		err = s.documentMetaService.PublishDocumentRagEvent(ctx, cmd)
		return err
	})
	return err
}

func (s *DocumentMetaAPI) SaveStatusBySourceType(ctx context.Context, cmd *command.DocumentMetaSaveStatusBySourceType) error {
	return s.documentMetaService.SaveStatusBySourceType(ctx, cmd)
}

func (s *DocumentMetaAPI) CreateMany(ctx context.Context, cmd *command.DocumentMetaCreateManyCommand) error {
	return s.documentMetaService.CreateMany(ctx, cmd.Data)
}

func (s *DocumentMetaAPI) Update(ctx context.Context, cmd *command.DocumentMetaUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"name", "value"})
	return s.documentMetaService.Update(ctx, &cmd.Data, opts)
}
