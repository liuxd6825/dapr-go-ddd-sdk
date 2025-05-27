package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	RegisterDrawio(app, baseUrl, env, rootPath)
}

func RegisterDrawio(app *iris.Application, baseUrl string, env *env.Env, rootPath string) {
	drawioAPI := NewDrawIoAPI(env, rootPath)
	mvc.Configure(app.Party(baseUrl), func(a *mvc.Application) {
		a.Handle(drawioAPI)
	})
}
