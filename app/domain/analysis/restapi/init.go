package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Register(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewSuTaskApi(env, baseUrl))
	restapi.RegisterController(app, NewSuTaskAccountApi(env, baseUrl))
	restapi.RegisterController(app, NewSuTaskRecordApi(env, baseUrl))
}
