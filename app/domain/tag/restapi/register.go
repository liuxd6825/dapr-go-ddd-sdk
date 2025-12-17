package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.tag"}, func() error {
		RegisterTagApi(app, baseUrl, env)
		RegisterTagTypeApi(app, baseUrl, env)
		RegisterTagRelationApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterTagApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewTagAPI(env, baseUrl))
}

func RegisterTagTypeApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewTagTypeAPI(env, baseUrl))
}

func RegisterTagRelationApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewTagRelationAPI(env, baseUrl))
}
