package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "draw"}, func() error {
		RegisterDrawGraphApi(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDrawGraphApi(app *iris.Application, baseUrl string, env *env.Env) {
	graphApi := NewGraphAPI(env, baseUrl)
	restapi.InitController(app, graphApi)
}
