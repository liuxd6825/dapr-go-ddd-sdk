package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type DocumentAPI struct {
	env        *env.Env
	docService *service.DocumentService
}

func NewDocumentAPI(env *env.Env, rootPath string) *DocumentAPI {
	docService := service.NewDocumentService()
	return &DocumentAPI{
		env:        env,
		docService: docService,
	}
}

func (s *DocumentAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/rag/document", "FindPaging")
	b.Handle(iris.MethodPost, "/rag/document", "Create")
	b.Handle(iris.MethodPut, "/rag/document", "Update")
	b.Handle(iris.MethodDelete, "/rag/document", "Delete")
	b.Handle(iris.MethodPost, "/rag/document:scan", "Scan")
}

func (s *DocumentAPI) FindPaging(ictx iris.Context) {

}

func (s *DocumentAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		return s.docService.Create(ctx, cmd)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *DocumentAPI) Update(ictx iris.Context) {

}

func (s *DocumentAPI) Delete(ictx iris.Context) {

}

func (s *DocumentAPI) Scan(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.DocumentScanCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.docService.Scan(ctx, cmd)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
