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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type ChatAPI struct {
	env         *env.Env
	chatService *service.ChatService
	msgService  *service.MessageService
}

func NewChatAPI(env *env.Env, rootPath string) *ChatAPI {
	chatService := service.NewChatService()
	msgService := service.NewMessageService()
	return &ChatAPI{
		env:         env,
		chatService: chatService,
		msgService:  msgService,
	}
}

func (s *ChatAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/chat", "Create")
	b.Handle(iris.MethodPut, "/rag/chat", "Update")
	b.Handle(iris.MethodPut, "/rag/chat:rename", "Rename")
	b.Handle(iris.MethodDelete, "/rag/chat", "Delete")
	b.Handle(iris.MethodGet, "/rag/chat", "FindPaging")
}

func (s *ChatAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatCreateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		vErr := errors.NewVerifyError()
		if cmd.Data.CaseId == "" {
			vErr.AppendField("caseId", "不能为空", "案件ID")
		}
		if vErr.HasError() {
			return vErr
		}
		s.chatService.Create(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *ChatAPI) Rename(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		opts := idao.NewCallOptions()
		opts.SetUpdateFields([]string{"title"})
		s.chatService.Update(ctx, &cmd.Data, opts)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *ChatAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatUpdateCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.chatService.Update(ctx, &cmd.Data)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *ChatAPI) Delete(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd *command.ChatDeleteCommand
		if err := ictx.ReadJSON(&cmd); err != nil {
			return err
		}
		s.msgService.DeleteByRSQL(ctx, "chat_id=='"+cmd.Data.Id+"'")
		s.chatService.DeleteById(ctx, cmd.Data.Id)
		return nil
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}

func (s *ChatAPI) FindPaging(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.URLParam("case-id")
		user, _ := appctx.GetAuthUser(ctx)
		qry := store.NewFindPagingQueryRequest()
		qry.PageNum = 0
		qry.PageSize = 99999999999999
		qry.Filter = "case_id=='" + caseId + "' and creator_id=='" + user.GetId() + "'"
		qry.Sort = "created_time:desc"
		qry.IsTotalRows = true
		res := s.chatService.FindPaging(ctx, qry)
		return web.SetData(ictx, res)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
