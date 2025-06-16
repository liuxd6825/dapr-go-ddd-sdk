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

type MessageAPI struct {
	env        *env.Env
	msgService *service.MessageService
}

func NewMessageAPI(env *env.Env, rootPath string) *MessageAPI {
	messageService := service.NewMessageService()
	return &MessageAPI{
		env:        env,
		msgService: messageService,
	}
}

func (s *MessageAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/message", "Create")
	b.Handle(iris.MethodPut, "/rag/message", "Update")
}

func (s *MessageAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.MessageCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.msgService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *MessageAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.MessageUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.msgService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
