package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	sysconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.card"}, func() error {

		RegisterHomeAPI(app, baseUrl, env, config)
		RegisterGroupAPI(app, baseUrl, env, config)
		RegisterCardAPI(app, baseUrl, env, config)
		RegisterAppAPI(app, baseUrl, env, config)
		RegisterFunAPI(app, baseUrl, env, config)
		RegisterCardFileAPI(app, baseUrl, env, config)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterHomeAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewHomeAPI(env, baseUrl, config.MongoDbKey))
}

func RegisterGroupAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewGroupAPI(env, baseUrl, config.MongoDbKey))
}

func RegisterCardAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewCardAPI(env, baseUrl, config.MongoDbKey))
}

func RegisterAppAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewAppAPI(env, baseUrl, config.MongoDbKey))
}

func RegisterFunAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewFunAPI(env, baseUrl, config.MongoDbKey))
}

func RegisterCardFileAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewCardFileAPI(env, baseUrl, config.MongoDbKey))
}
