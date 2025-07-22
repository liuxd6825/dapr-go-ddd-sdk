package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "import"}, func() error {
		RegisterExcelRowApi(app, baseUrl, env)
		RegisterExcelFileApi(app, baseUrl, env)
		RegisterExcelSheetApi(app, baseUrl, env)
		RegisterTaskApi(app, baseUrl, env)
		RegisterRecordAPI(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterExcelRowApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewExcelRowApi(env, baseUrl)
	restapi.InitController(app, api)
}

func RegisterExcelFileApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewExcelFileApi(env, baseUrl)
	restapi.InitController(app, api)
}

func RegisterExcelSheetApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewExcelSheetApi(env, baseUrl)
	restapi.InitController(app, api)
}

func RegisterTaskApi(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewTaskAPI(env, baseUrl)
	restapi.InitController(app, api)
}

func RegisterRecordAPI(app *iris.Application, baseUrl string, env *env.Env) {
	api := NewRecordAPI(env, baseUrl)
	restapi.InitController(app, api)
}
