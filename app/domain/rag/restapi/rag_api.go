package restapi

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RagAPI struct {
	env        *env.Env
	ragService *service.RagService
	rootPath   string
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
	return &RagAPI{
		env:        env,
		ragService: ragService,
		rootPath:   rootPath,
	}
}

func (s *RagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/query", "Query")
	b.Handle(iris.MethodPost, "/rag/create-tenant", "CreateTenant")
	b.Handle(iris.MethodPost, "/rag/create-case", "CreateCase")
}

func (s *RagAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath+"/rag", s)
	ctl.Post("query", "Query")
	ctl.Post("create-tenant", "CreateTenant")
	ctl.Post("create-case", "CreateCase")
	return nil
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
