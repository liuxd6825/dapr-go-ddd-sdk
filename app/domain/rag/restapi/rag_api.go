package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type RagAPI struct {
	env          *env.Env
	graphService *service.GraphService
}

func NewRagAPI(env *env.Env, rootPath string) *RagAPI {
	graphService := service.NewGraphService()
	return &RagAPI{
		env:          env,
		graphService: graphService,
	}
}

func (s *RagAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/case/{caseId}/rag", "Query")
}

func (s *RagAPI) Query(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.Params().Get("caseId")
		query := ictx.URLParamDefault("query", "")
		if query == "" {
			return errors.New("query is required")
		}
		graphView := s.graphService.FindInCaseByNames(ctx, caseId, names)
		result := response.NewResultList(graphView)
		return web.SetData(ictx, result)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
