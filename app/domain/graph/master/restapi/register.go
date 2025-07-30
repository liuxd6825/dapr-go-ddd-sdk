package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	draw_restapi "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	_ = logs.DebugStart(context.Background(), logs.Fields{"service name ": "graph"}, func() error {
		draw_restapi.RegisterAllApi(app, baseUrl, env)
		return nil
	})
}
