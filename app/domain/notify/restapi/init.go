package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	RegisterNotifyApi(app, baseUrl, env)
}

func RegisterNotifyApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewNotifyApi(baseUrl, env))
}
