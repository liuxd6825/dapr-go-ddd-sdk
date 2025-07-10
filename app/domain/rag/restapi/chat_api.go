package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type ChatAPI struct {
	env         *env.Env
	chatService *service.ChatService
	msgService  *service.MessageService
	rootPath    string
}

func NewChatAPI(env *env.Env, rootPath string) *ChatAPI {
	chatService := service.NewChatService()
	msgService := service.NewMessageService()
	return &ChatAPI{
		env:         env,
		chatService: chatService,
		msgService:  msgService,
		rootPath:    rootPath,
	}
}

func (s *ChatAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/chat", "Create")
	b.Handle(iris.MethodPut, "/rag/chat", "Update")
	b.Handle(iris.MethodPut, "/rag/chat:rename", "Rename")
	b.Handle(iris.MethodDelete, "/rag/chat", "Delete")
	b.Handle(iris.MethodGet, "/rag/chat", "FindPaging")
}

func (s *ChatAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath+"/rag", s)
	ctl.Post("chat", "Create")
	ctl.Put("chat", "Update")
	ctl.Put("chat:rename", "Rename")
	ctl.Delete("chat", "Delete")
	ctl.GetPaging("chat", "FindPaging")
	return nil
}

func (s *ChatAPI) Create(ctx context.Context, cmd *command.ChatCreateCommand) error {
	vErr := errors.NewVerifyError()
	if cmd.Data.CaseId == "" {
		vErr.AppendField("caseId", "不能为空", "案件ID")
	}
	if vErr.HasError() {
		return vErr
	}
	return s.chatService.Create(ctx, &cmd.Data).GetError()
}

func (s *ChatAPI) Rename(ctx context.Context, cmd *command.ChatUpdateCommand) error {
	opts := idao.NewCallOptions().SetUpdateFields([]string{"title"})
	return s.chatService.Update(ctx, &cmd.Data, opts).GetError()
}

func (s *ChatAPI) Update(ctx context.Context, cmd *command.ChatUpdateCommand) error {
	return s.chatService.Update(ctx, &cmd.Data).GetError()
}

func (s *ChatAPI) Delete(ctx context.Context, cmd *command.ChatDeleteCommand) {
	s.msgService.DeleteByRSQL(ctx, "chat_id=='"+cmd.Data.Id+"'")
	s.chatService.DeleteById(ctx, cmd.Data.Id)
}

func (s *ChatAPI) FindPaging(ctx context.Context, ictx iris.Context, qry *store.FindPagingQueryRequest) (any, error) {
	caseId := ictx.URLParam("case-id")
	user, _ := appctx.GetAuthUser(ctx)
	qry.PageNum = 0
	qry.PageSize = 99999999999999
	qry.Filter = "case_id=='" + caseId + "' and creator_id=='" + user.GetId() + "'"
	qry.Sort = "created_time:desc"
	qry.IsTotalRows = true
	res := s.chatService.FindPaging(ctx, qry)
	return res, res.GetError()
}
