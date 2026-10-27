package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	drawconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name ": "draw"}, func() error {
		drawconfigs.Init(env)
		if drawconfigs.GetEnable() {
			RegisterDrawApi(app, baseUrl, env)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDrawApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewDrawIoAPI(env, baseUrl, drawconfigs.GetConfig().HugeDbKey))
}
