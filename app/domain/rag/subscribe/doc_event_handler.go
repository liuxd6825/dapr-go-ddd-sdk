package subscribe

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/outside"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
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
	ctl.EventHandle("document-create-document", "CreateDocumentEvent")
	ctl.Handle(iris.MethodOptions, "document-create-document", "Check")
	return ctl
}

func (s *DocumentEventHandler) Check(ctx context.Context) error {
	return nil
}

func (s *DocumentEventHandler) CreateDocumentEvent(ctx context.Context, event *event.DocumentCreateEvent) error {
	logs.DebugEvent(ctx, event, "CreateDocumentEvent")
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
	cmd.Data.FsKey = event.Data.FsKey
	cmd.Data.CaseId = event.Data.CaseId
	cmd.Data.FilePath = event.Data.FilePath
	cmd.Data.FileId = event.Data.FileId
	cmd.Data.FileName = event.Data.FileName
	return cmd, nil
}
