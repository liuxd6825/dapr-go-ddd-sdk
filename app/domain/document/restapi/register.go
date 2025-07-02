package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "document"}, func() error {
		RegisterFolderApi(app, baseUrl, env, rootPath)
		RegisterFileApi(app, baseUrl, env, rootPath)
		RegisterDocumentApi(app, baseUrl, env, rootPath)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterFolderApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	ragApi := NewFolderAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(ragApi)
	})
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewDocumentAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterFileApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewFileAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}
