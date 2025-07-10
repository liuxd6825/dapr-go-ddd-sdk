package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	RegisterRagApi(app, baseUrl, env)
	RegisterChatApi(app, baseUrl, env)
	RegisterMessageApi(app, baseUrl, env)
	RegisterDocumentApi(app, baseUrl, env)
}

func RegisterRagApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.InitController(app, NewRagAPI(env, baseUrl))
}

func RegisterMessageApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.InitController(app, NewMessageAPI(env, baseUrl))
}

func RegisterChatApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.InitController(app, NewChatAPI(env, baseUrl))
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.InitController(app, NewDocumentAPI(env, baseUrl))
}
