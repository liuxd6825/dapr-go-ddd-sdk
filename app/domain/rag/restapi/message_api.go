package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
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

func (s *MessageAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/rag", "rag.MessageAPI", s)
	ctl.Post("message", "Create")
	ctl.Put("message", "Update")
	ctl.GetPaging("/{chatId}/messages", "GetByChatId")
	return ctl
}

func (s *MessageAPI) Create(ctx context.Context, cmd *command.MessageCreateCommand) error {
	return s.msgService.Create(ctx, &cmd.Data)
}

func (s *MessageAPI) Update(ctx context.Context, cmd *command.MessageUpdateCommand) error {
	return s.msgService.Update(ctx, &cmd.Data).GetError()
}

func (s *MessageAPI) GetByChatId(ctx context.Context, ictx iris.Context, qry *restapi.FindPagingRequest) (any, error) {
	chatId := ictx.Params().GetString("chatId")
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.MustFilter = "chat_id=='" + chatId + "'"
	qry.Sort = "order_num:asc"
	qry.IsTotalRows = true
	res := s.msgService.FindPaging(ctx, qry)
	return res, res.GetError()
}
