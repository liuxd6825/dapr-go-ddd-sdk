package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"strings"
)

type GraphAPI struct {
	env          *env.Env
	graphService *service.GraphService
}

func NewGraphAPI(env *env.Env, rootPath string) *GraphAPI {
	graphService := service.NewGraphService()
	return &GraphAPI{
		env:          env,
		graphService: graphService,
	}
}

func (s *GraphAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/case/{caseId}/graph?name={name}", "FindByName")
}

func (s *GraphAPI) FindByName(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		caseId := ictx.Params().Get("caseId")
		nameStr := ictx.Params().Get("name")
		names := strings.Split(nameStr, ",")
		graphView := s.graphService.FindInCaseByNames(ctx, caseId, names)
		result := response.NewResultList(graphView)
		return web.SetData(ictx, result)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
