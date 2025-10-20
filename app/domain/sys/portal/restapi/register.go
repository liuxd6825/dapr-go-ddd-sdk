package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.portal"}, func() error {
		RegisterTenantAPI(app, baseUrl, env)
		RegisterUserAPI(app, baseUrl, env)
		RegisterAuthAPI(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterTenantAPI(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewTenantAPI(env, baseUrl))
}

func RegisterUserAPI(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewUserAPI(env, baseUrl))
}

func RegisterAuthAPI(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewAuthAPI(env, baseUrl))
}
