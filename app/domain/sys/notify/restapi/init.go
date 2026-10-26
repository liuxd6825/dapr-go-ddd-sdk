package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	RegisterNotifyApi(app, baseUrl, env)
}

func RegisterNotifyApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewNotifyApi(baseUrl, env))
}
