package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
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
	b.Handle(iris.MethodGet, "/draw/{drawId}/graph", "FindByDrawId")
}

func (s *GraphAPI) FindByDrawId(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().Get("drawId")
		graphView := s.graphService.FindById(ctx, id)
		result := response.NewResultWithGraphView(graphView)
		return web.SetData(ictx, result)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
