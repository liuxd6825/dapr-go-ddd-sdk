package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "master"}, func() error {
		RegisterSchema(app, baseUrl, env, rootPath)
		RegisterHtml(app, baseUrl, env, rootPath)
		RegisterCdcToNeo4j(app, "", env, rootPath)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterSchema(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	schemaAPI := NewSchemaAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(schemaAPI)
	})
}

func RegisterHtml(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	htmlAPI := NewHtmlAPI(env, "web")
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(htmlAPI)
	})
}

func RegisterCdcToNeo4j(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	graphAPI := NewGraphAPI(env, rootPath)
	cdcAPI := NewCdcAPI(env, "")
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(graphAPI)
	})
	mvc.Configure(app.Party(""), func(a *mvc.Application) {
		a.Handle(cdcAPI)
	})
}
