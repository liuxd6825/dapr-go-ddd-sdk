package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	RegisterRagApi(app, baseUrl, env, rootPath)
	RegisterChatApi(app, baseUrl, env, rootPath)
	RegisterMessageApi(app, baseUrl, env, rootPath)
	RegisterDocumentApi(app, baseUrl, env, rootPath)
}

func RegisterRagApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	ragApi := NewRagAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(ragApi)
	})
}

func RegisterMessageApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewMessageAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterChatApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewChatAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	api := NewDocumentAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(api)
	})
}
