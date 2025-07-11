package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "tag"}, func() error {
		RegisterTagApi(app, baseUrl, env, rootPath)
		RegisterTagTypeApi(app, baseUrl, env, rootPath)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterTagApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	ragApi := NewTagAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(ragApi)
	})
}

func RegisterTagTypeApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewTagTypeAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}
