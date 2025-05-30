package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/graph/vis_network"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type GraphAPI struct {
	env          *env.Env
	graphService *service.GraphService
}

func NewGraphAPI(env *env.Env, rootPath string) *DrawAPI {
	graphService := service.NewGraphService()
	return &DrawAPI{
		env:          env,
		graphService: graphService,
	}
}

func (s *GraphAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/draw/graph/{draw-id}", "FindByDrawId")
}

func (s *GraphAPI) FindByDrawId(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().Get("draw-id")
		graphView := s.graphService.FindById(ctx, id)
		result := vis_network.NewResultWithGraphView(graphView)
		return web.SetData(ictx, result)
	}).Catch(func(ctx context.Context, err error) {
		web.SetError(ictx, err)
	})
}
