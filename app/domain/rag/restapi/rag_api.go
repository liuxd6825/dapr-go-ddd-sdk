package restapi

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type RagAPI struct {
	env        *env.Env
	ragService *service.RagService
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
	}
}

func (s *RagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/rag/query", "Query")
}

func (s *RagAPI) Query(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var query *my_rag.QueryParam
		if err := ictx.ReadJSON(&query); err != nil {
			return err
		}
		vErr := errors.NewVerifyError()
		if query.Query == "" {
			vErr.AppendField("query", "不能为空", "查询内容")
		}
		if query.CaseId == "" {
			vErr.AppendField("caseId", "不能为空", "项目ID")
		}
		if vErr.HasError() {
			return vErr
		}

		isPrint := logs.GetLevel() >= logs.InfoLevel
		ictx.Header("Content-Type", "text/event-stream")
		_, err := s.ragService.Query(ctx, *query, func(txt string) {
			_, _ = ictx.Writef(txt)
			if isPrint {
				fmt.Print(txt)
			}
			ictx.ResponseWriter().Flush()
		})
		if isPrint {
			fmt.Print("\n")
		}
		return err
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
