package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	RegisterRagApi(app, baseUrl, env)
	RegisterChatApi(app, baseUrl, env)
	RegisterMessageApi(app, baseUrl, env)
	RegisterDocumentApi(app, baseUrl, env)
}

func RegisterRagApi(app *iris.Application, baseUrl string, env *env.Env) {
	ragApi := NewRagAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(ragApi)
	})
}

func RegisterMessageApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewMessageAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterChatApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewChatAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewDocumentAPI(env, baseUrl)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}
