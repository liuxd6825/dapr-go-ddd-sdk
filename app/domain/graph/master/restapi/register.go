package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	graphconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env) error {
	return logs.DebugStart(context.Background(), logs.Fields{"service name ": "graph"}, func() error {
		graphconfigs.Init(env)
		RegisterDrawApi(app, baseUrl, env)
		RegisterRecordParquetApi(app, baseUrl, env, graphconfigs.GetConfig())
		//RegisterRecordNebulaApi(app, baseUrl, env)
		return nil
	})

}

func RegisterDrawApi(app *iris.Application, baseUrl string, env *envs.Env) {
	restapi.RegisterController(app, NewMasterAPI(env, baseUrl))
}

// RegisterRecordParquetApi 注册银行流水 parquet 导入 API
func RegisterRecordParquetApi(app *iris.Application, baseUrl string, env *envs.Env, cfg graphconfigs.Config) {
	recordService := service.NewRecordServiceWithEnv(env, cfg)
	restapi.RegisterController(app, NewRecordParquetAPI(env, recordService, baseUrl))
}

/*func RegisterRecordNebulaApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewRecordNebulaAPI(env, baseUrl))
}
*/
