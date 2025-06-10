package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type RagAPI struct {
	env        *env.Env
	ragService *service.RagService
}

func NewRagAPI(env *env.Env, rootPath string) *RagAPI {
	ragService := service.NewRagService()
	return &RagAPI{
		env:        env,
		ragService: ragService,
	}
}

func (s *RagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/case/{caseId}/rag", "Query")
}

func (s *RagAPI) Query(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.Params().Get("caseId")
		queryVal := ictx.URLParamDefault("query", "")
		if queryVal == "" {
			return errors.New("query is required")
		}
		deepVal := ictx.URLParamIntDefault("deep", 5)
		q := my_rag.QueryParam{
			CaseId:  caseId,
			Query:   queryVal,
			MaxDeep: deepVal,
		}
		ictx.Header("Content-Type", "text/event-stream")
		_, err := s.ragService.Query(ctx, q, func(txt string) {
			ictx.Writef(txt)
			ictx.ResponseWriter().Flush()
		})
		return err
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
