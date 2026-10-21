package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	sysconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env, cfg sysconfigs.Config) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.currency"}, func() error {
		RegisterDictAPI(app, baseUrl, env, cfg.MongoDbKey)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDictAPI(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCurrencyAPI(env, baseUrl, dbKey))
}
