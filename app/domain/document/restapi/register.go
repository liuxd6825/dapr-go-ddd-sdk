package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "document"}, func() error {
		RegisterFolderApi(app, baseUrl, env)
		RegisterFileApi(app, baseUrl, env)
		RegisterDocumentApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterFolderApi(app *iris.Application, baseUrl string, env *env.Env) {
	ragApi := NewFolderAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(ragApi)
	})
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewDocumentAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterFileApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewFileAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}
