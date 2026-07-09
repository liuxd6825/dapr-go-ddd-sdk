package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/outside"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DocumentEventHandler struct {
	rootPath   string
	env        *env.Env
	docService *service.DocumentService
}

func NewDocumentEventHandler(env *env.Env, baseUrl string) *DocumentEventHandler {
	return &DocumentEventHandler{
		rootPath:   baseUrl,
		env:        env,
		docService: service.NewDocumentService(outside.NewDocStatusProvider()),
	}
}

func (s *DocumentEventHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, "subscribe/rag/event", "RagDocumentEventSubHandler", s)
	ctl.EventHandle("document-create-event", "CreateDocumentEvent")
	ctl.EventHandle("document-delete-event", "DeleteDocumentEvent")
	ctl.Handle(iris.MethodOptions, "document-create-event", "Check")
	return ctl
}

func (s *DocumentEventHandler) Check(ctx context.Context) error {
	return nil
}

func (s *DocumentEventHandler) CreateDocumentEvent(ctx context.Context, event *event.DocumentCreateEvent) error {
	fmt.Printf("%v=====>rag-CreateDocumentEvent=====>%s", time.Now(), event)

	cmd, err := newDocumentCreateCommand(event)
	if err != nil {
		return err
	}
	err = s.docService.Create(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func newDocumentCreateCommand(event *event.DocumentCreateEvent) (*command.DocumentCreateCommand, error) {
	cmd := &command.DocumentCreateCommand{}
	cmd.CommandId = event.Id
	cmd.Data.Id = event.Id
	cmd.Data.FsKey = event.Data.FsKey
	cmd.Data.CaseId = event.Data.CaseId
	cmd.Data.FilePath = event.Data.FilePath
	cmd.Data.FileId = event.Data.FileId
	cmd.Data.FileName = event.Data.FileName
	cmd.Data.SourceType = event.Data.SourceType
	cmd.Data.SourceId = event.Data.SourceId
	cmd.Data.SourceApp = event.Data.SourceApp
	cmd.Data.SourceUrl = event.Data.SourceUrl
	cmd.Data.SourceName = event.Data.SourceName
	return cmd, nil
}

func (s *DocumentEventHandler) DeleteDocumentEvent(ctx context.Context, event *event.DocumentDeleteEvent) error {
	fmt.Printf("%v=====>rag-DeleteDocumentEvent=====>%s", time.Now(), event)

	cmd, err := newDocumentDeleteCommand(event)
	if err != nil {
		return err
	}
	err = s.docService.DeleteByRSql(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
func newDocumentDeleteCommand(event *event.DocumentDeleteEvent) (*command.DocumentDeleteByDocCommand, error) {
	cmd := &command.DocumentDeleteByDocCommand{}
	cmd.CommandId = event.Id
	cmd.Data.DocId = event.Data.DocId
	return cmd, nil
}
