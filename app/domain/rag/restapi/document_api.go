package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/outside"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentAPI struct {
	env        *env.Env
	docService *service.DocumentService
	rootPath   string
}

func NewDocumentAPI(env *env.Env, rootPath string) *DocumentAPI {
	docService := service.NewDocumentService(outside.NewImportStatusProvider())
	return &DocumentAPI{
		env:        env,
		docService: docService,
		rootPath:   rootPath,
	}
}

func (s *DocumentAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/rag/document", "FindPaging")
	b.Handle(iris.MethodPost, "/rag/document", "Create")
	b.Handle(iris.MethodPut, "/rag/document", "Update")
	b.Handle(iris.MethodDelete, "/rag/document", "Delete")
	b.Handle(iris.MethodPost, "/rag/document:scan", "Scan")
}

func (s *DocumentAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/rag", "rag.DocumentAPI", s)
	ctl.GetPaging("document", "FindPaging")
	ctl.Post("document", "Create")
	ctl.Put("document", "Update")
	ctl.Delete("document", "Delete")
	ctl.Post("document:scan", "Scan")
	return ctl
}

func (s *DocumentAPI) FindPaging(ictx iris.Context) {

}

func (s *DocumentAPI) Create(ctx context.Context, cmd *command.DocumentCreateCommand) error {
	return s.docService.Create(ctx, cmd)
}

func (s *DocumentAPI) Update(ictx iris.Context) {

}

func (s *DocumentAPI) Delete(ictx iris.Context) {

}

func (s *DocumentAPI) Scan(ctx context.Context, cmd *command.DocumentScanCommand) {
	s.docService.Scan(ctx, cmd)
}
