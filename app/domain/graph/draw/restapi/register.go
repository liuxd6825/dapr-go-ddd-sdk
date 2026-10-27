package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	graphconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "draw"}, func() error {
		graphconfigs.Init(env)
		RegisterDrawGraphApi(app, baseUrl, env, graphconfigs.GetGraphDbKey())
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDrawGraphApi(app *iris.Application, baseUrl string, env *envs.Env, graphDbKey string) {
	graphApi := NewGraphAPI(env, baseUrl, graphDbKey)
	restapi.RegisterController(app, graphApi)
}
