package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RagAPI struct {
	env               *env.Env
	ragService        *service.RagService
	chatService       *service.ChatService
	messageService    *service.MessageService
	chatSummaryService *service.ChatSummaryService
	rootPath          string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type QueryRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float32  `json:"temperature"`
}

func NewRagAPI(env *env.Env, rootPath string) *RagAPI {
	ragService := service.NewRagService()
	chatService := service.NewChatService()
	compressConfig := config.DefaultCompressConfig()
	chatSummaryService := service.NewChatSummaryService(compressConfig)
	messageService := service.NewMessageService()

	messageService.SetChatSummaryService(chatSummaryService)
	messageService.SetCompressConfig(compressConfig)
	ragService.SetMessageService(messageService)
	ragService.SetChatSummaryService(chatSummaryService)

	return &RagAPI{
		env:               env,
		ragService:        ragService,
		chatService:       chatService,
		messageService:    messageService,
		chatSummaryService: chatSummaryService,
		rootPath:          rootPath,
	}
}

func (s *RagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/query", "Query")
	b.Handle(iris.MethodPost, "/rag/create-tenant", "CreateTenant")
	b.Handle(iris.MethodPost, "/rag/create-case", "CreateCase")
	b.Handle(iris.MethodPost, "/rag/compress", "Compress")
	b.Handle(iris.MethodPost, "/rag/force-unlock", "ForceUnlock")
	b.Handle(iris.MethodPost, "/rag/is-locked", "IsLocked")
}

func (s *RagAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/rag", "rag.RagAPI", s)
	ctl.Post("query", "Query")
	ctl.Post("create-tenant", "CreateTenant")
	ctl.Post("create-case", "CreateCase")
	return ctl
}

// CreateTenant 加载租户知识库数据
func (s *RagAPI) CreateTenant(ctx context.Context) error {
	return s.ragService.CreateTenant(ctx)
}

// CreateCase 加载租户知识库数据
func (s *RagAPI) CreateCase(ctx context.Context, cmd *command.RagCreateCaseCommand) error {
	return s.ragService.CreateCase(ctx, cmd.Data.CaseId)
}

// Query 查询
func (s *RagAPI) Query(ctx context.Context, ictx iris.Context, query *my_rag.QueryParam) error {
	isPrint := logs.GetLevel() >= logs.InfoLevel
	ictx.Header("Content-Type", "text/event-stream")
	_, err := s.ragService.Query(ctx, *query, func(txt string) {
		_, _ = ictx.Writef(txt)
		if isPrint {
			fmt.Print(txt)
		}
		ictx.ResponseWriter().Flush()
	})
	if err == nil && isPrint {
		//fmt.Print("\n")
	}
	return err
}

type CompressRequest struct {
	ChatId string `json:"chatId"`
	Force  bool   `json:"force"`
}

type CompressResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	TokensBefore int    `json:"tokensBefore"`
	TokensAfter  int    `json:"tokensAfter"`
}

type LockResponse struct {
	Locked bool   `json:"locked"`
	ChatId string `json:"chatId"`
}

func (s *RagAPI) Compress(ctx context.Context, req *CompressRequest) (*CompressResponse, error) {
	locked := s.messageService.IsLocked(ctx, req.ChatId)
	if locked {
		return &CompressResponse{
			Success: false,
			Message: "chat is locked, please try again later",
		}, nil
	}

	s.chatSummaryService.AcquireLock(ctx, req.ChatId)
	defer s.chatSummaryService.ReleaseLock(ctx, req.ChatId)

	result, err := s.chatSummaryService.Compress(ctx, req.ChatId, req.Force, nil)
	if err != nil {
		return &CompressResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &CompressResponse{
		Success:      true,
		Message:      "compress completed",
		TokensBefore: result.OriginalTokens,
		TokensAfter:  result.CompressedTokens,
	}, nil
}

func (s *RagAPI) ForceUnlock(ctx context.Context, req *CompressRequest) error {
	return s.chatSummaryService.ForceUnlock(ctx, req.ChatId)
}

func (s *RagAPI) IsLocked(ctx context.Context, req *CompressRequest) (*LockResponse, error) {
	locked := s.messageService.IsLocked(ctx, req.ChatId)
	return &LockResponse{
		Locked: locked,
		ChatId: req.ChatId,
	}, nil
}
