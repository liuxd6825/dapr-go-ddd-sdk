package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	_ = logs.DebugStart(context.Background(), logs.Fields{"service name ": "graph"}, func() error {
		RegisterGraphApi(app, baseUrl, env, rootPath)
		return nil
	})
}

func RegisterGraphApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	graphApi := NewGraphAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(graphApi)
	})
}
