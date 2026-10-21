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
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.portal"}, func() error {
		sysconfigs.Init(env)
		mongoDbKey := config.MongoDbKey
		RegisterTenantAPI(app, baseUrl, env, mongoDbKey)
		RegisterUserAPI(app, baseUrl, env, mongoDbKey)
		RegisterAuthAPI(app, baseUrl, env, mongoDbKey)
		RegisterTenantUserAPI(app, baseUrl, env, mongoDbKey)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterTenantAPI(app *iris.Application, baseUrl string, env *envs.Env, mongoDbKey string) {
	restapi.RegisterController(app, NewTenantAPI(env, baseUrl, mongoDbKey))
}

func RegisterUserAPI(app *iris.Application, baseUrl string, env *envs.Env, mongoDbKey string) {
	restapi.RegisterController(app, NewUserAPI(env, baseUrl, mongoDbKey))
}

func RegisterAuthAPI(app *iris.Application, baseUrl string, env *envs.Env, mongoDbKey string) {
	restapi.RegisterController(app, NewAuthAPI(env, baseUrl, mongoDbKey))
}

func RegisterTenantUserAPI(app *iris.Application, baseUrl string, env *envs.Env, mongoDbKey string) {
	restapi.RegisterController(app, NewTenantUserAPI(env, baseUrl, mongoDbKey))
}
