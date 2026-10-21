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
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "sys.case"}, func() error {
		RegisterDictAPI(app, baseUrl, env, config)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RegisterDictAPI(app *iris.Application, baseUrl string, env *envs.Env, config sysconfigs.Config) {
	restapi.RegisterController(app, NewCaseAPI(env, baseUrl, config.MongoDbKey))
}
