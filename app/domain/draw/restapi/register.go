package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "draw"}, func() error {
		RegisterDrawApi(app, baseUrl, env, rootPath)
		RegisterGraphApi(app, baseUrl, env, rootPath)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDrawApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	drawioAPI := NewDrawIoAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(drawioAPI)
	})
}
func RegisterGraphApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	graphApi := NewGraphAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(graphApi)
	})
}
