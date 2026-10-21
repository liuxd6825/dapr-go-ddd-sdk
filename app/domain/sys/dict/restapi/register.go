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
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.dictType"}, func() error {
		RegisterDictTypeAPI(app, baseUrl, env, config.MongoDbKey)
		RegisterDictAPI(app, baseUrl, env, config.MongoDbKey)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDictTypeAPI(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewDictTypeAPI(env, baseUrl, dbKey))
}

func RegisterDictAPI(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewDictAPI(env, baseUrl, dbKey))
}
