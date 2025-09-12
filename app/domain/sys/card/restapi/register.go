package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.card"}, func() error {
		RegisterHomeAPI(app, baseUrl, env)
		RegisterGroupAPI(app, baseUrl, env)
		RegisterCardAPI(app, baseUrl, env)
		RegisterAppAPI(app, baseUrl, env)
		RegisterFunAPI(app, baseUrl, env)
		RegisterCardFileAPI(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterHomeAPI(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewHomeAPI(env, baseUrl))
}

func RegisterGroupAPI(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewGroupAPI(env, baseUrl))
}

func RegisterCardAPI(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewCardAPI(env, baseUrl)
	restapi.RegisterController(app, api)
}

func RegisterAppAPI(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewAppAPI(env, baseUrl)
	restapi.RegisterController(app, api)
}

func RegisterFunAPI(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewFunAPI(env, baseUrl)
	restapi.RegisterController(app, api)
}

func RegisterCardFileAPI(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewCardFileAPI(env, baseUrl)
	restapi.RegisterController(app, api)
}
