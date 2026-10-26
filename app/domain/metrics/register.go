package metrics

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "metrics"}, func() error {
		RegisterMetricsAPI(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterMetricsAPI(app *iris.Application, baseUrl string, env *envs.Env) {
	app.Get("/metrics", func(ctx iris.Context) {
		promhttp.Handler().ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	})
}
