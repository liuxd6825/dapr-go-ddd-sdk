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

type FolderMetaAPI struct {
	rootPath          string
	env               *env.Env
	folderMetaService *service.FolderMetaService
}

func NewFolderMetaAPI(env *env.Env, rootPath string) *FolderMetaAPI {
	folderMetaService := service.NewFolderMetaService()
	return &FolderMetaAPI{
		rootPath:          rootPath,
		env:               env,
		folderMetaService: folderMetaService,
	}
}

func (s *FolderMetaAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/doc", "folder.FolderMetaAPI", s)
	ctl.Post("/folder-meta", "Create")
	ctl.Post("/folder-meta:batch", "CreateMany")
	ctl.Put("/folder-meta", "Update")
	return ctl
}

func (s *FolderMetaAPI) Create(ctx context.Context, cmd *command.FolderMetaCreateCommand) error {
	err := s.folderMetaService.Create(ctx, &cmd.Data)
	return err
}

func (s *FolderMetaAPI) CreateMany(ctx context.Context, cmd *command.FolderMetaCreateManyCommand) error {
	return s.folderMetaService.CreateMany(ctx, cmd.Data)
}

func (s *FolderMetaAPI) Update(ctx context.Context, cmd *command.FolderMetaUpdateCommand) error {
	opts := idao.NewCallOptions()
	opts.SetUpdateFields([]string{"name", "value"})
	return s.folderMetaService.Update(ctx, &cmd.Data, opts)
}
