package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentMetaAPI struct {
	rootPath            string
	env                 *env.Env
	documentMetaService *service.DocumentMetaService
	docService          *service.DocumentService
}

func NewDocumentMetaAPI(env *env.Env, rootPath string) *DocumentMetaAPI {
	documentMetaService := service.NewDocumentMetaService()
	docService := service.NewDocumentService()
	return &DocumentMetaAPI{
		rootPath:            rootPath,
		env:                 env,
		documentMetaService: documentMetaService,
		docService:          docService,
	}
}

func (s *DocumentMetaAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "document.DocumentMetaAPI", s)
	ctl.Post("/document-meta", "Create")
	ctl.Post("/document-meta:submit", "Submit")
	ctl.Post("/document-meta:batch", "CreateMany")
	ctl.Put("/document-meta", "Update")
	return ctl
}

func (s *DocumentMetaAPI) Create(ctx context.Context, cmd *command.DocumentMetaCreateCommand) error {
	err := s.documentMetaService.Create(ctx, &cmd.Data)
	return err
}

func (s *DocumentMetaAPI) Submit(ctx context.Context, cmd *command.DocumentMetaCreateCommand) error {
	meta, err := s.documentMetaService.FindByDocumentIdAndName(ctx, cmd.Data.DocumentId, cmd.Data.Name)
	if err != nil {
		return err
	}

	if meta != nil {
		meta.Value = cmd.Data.Value
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"value"})
		err = s.documentMetaService.Update(ctx, meta, opts)
	} else {
		err = s.documentMetaService.Create(ctx, &cmd.Data)
	}

	return err
}

func (s *DocumentMetaAPI) CreateMany(ctx context.Context, cmd *command.DocumentMetaCreateManyCommand) error {
	return s.documentMetaService.CreateMany(ctx, cmd.Data)
}

func (s *DocumentMetaAPI) Update(ctx context.Context, cmd *command.DocumentMetaUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"name", "value"})
	return s.documentMetaService.Update(ctx, &cmd.Data, opts)
}
