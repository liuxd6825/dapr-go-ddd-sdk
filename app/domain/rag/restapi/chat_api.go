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

type ChatAPI struct {
	env         *env.Env
	chatService *service.ChatService
}

func NewChatAPI(env *env.Env, rootPath string) *ChatAPI {
	chatService := service.NewChatService()
	return &ChatAPI{
		env:         env,
		chatService: chatService,
	}
}

func (s *ChatAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/chat", "Create")
	b.Handle(iris.MethodPut, "/rag/chat", "Update")
}

func (s *ChatAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.chatService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *ChatAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.chatService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
