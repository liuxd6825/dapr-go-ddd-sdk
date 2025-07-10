package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type MessageAPI struct {
	env        *env.Env
	msgService *service.MessageService
	rootPath   string
}

func NewMessageAPI(env *env.Env, rootPath string) *MessageAPI {
	messageService := service.NewMessageService()
	return &MessageAPI{
		env:        env,
		msgService: messageService,
		rootPath:   rootPath,
	}
}

func (s *MessageAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/message", "Create")
	b.Handle(iris.MethodPut, "/rag/message", "Update")
	b.Handle(iris.MethodGet, "/rag/{chatId}/messages", "GetByChatId")
}

func (s *MessageAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath+"/rag", s)
	ctl.Post("message", "Create")
	ctl.Put("message", "Update")
	ctl.GetOne("/{chatId}/messages", "GetByChatId")
	return nil
}

func (s *MessageAPI) Create(ctx context.Context, cmd *command.MessageCreateCommand) error {
	return s.msgService.Create(ctx, &cmd.Data).GetError()
}

func (s *MessageAPI) Update(ctx context.Context, cmd *command.MessageUpdateCommand) error {
	return s.msgService.Update(ctx, &cmd.Data).GetError()
}

func (s *MessageAPI) GetByChatId(ictx iris.Context) {
	restapi.Try(ictx, func(ctx context.Context) error {
		chatId := ictx.Params().GetString("chatId")
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "chat_id=='" + chatId + "'"
		qry.Sort = "order_num:asc"
		qry.IsTotalRows = true
		res := s.msgService.FindPaging(ctx, qry)
		return restapi.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		restapi.SetError(ictx, err)
	})
}
