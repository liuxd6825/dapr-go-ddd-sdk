package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "master"}, func() error {
		RegisterSchema(app, baseUrl, env)
		RegisterHtml(app, baseUrl, env)
		RegisterCdcToNeo4j(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterSchema(app *iris.Application, baseUrl string, env *env.Env) {
	schemaAPI := NewSchemaAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(schemaAPI)
	})
}

func RegisterHtml(app *iris.Application, baseUrl string, env *env.Env) {
	htmlAPI := NewHtmlAPI(env, "web")
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(htmlAPI)
	})
}

func RegisterCdcToNeo4j(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.InitController(app, NewGraphAPI(env, baseUrl))
	restapi.InitController(app, NewCdcAPI(env, baseUrl))
}
