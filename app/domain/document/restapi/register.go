package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "document"}, func() error {
		RegisterFolderApi(app, baseUrl, env)
		RegisterFileApi(app, baseUrl, env)
		RegisterDocumentApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterFolderApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewFolderAPI(env, baseUrl))
}

func RegisterDocumentApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewDocumentAPI(env, baseUrl))
}

func RegisterFileApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewFileAPI(env, baseUrl))
}
