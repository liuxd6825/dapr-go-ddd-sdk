package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	_ = logs.DebugStart(context.Background(), logs.Fields{"service name ": "graph"}, func() error {
		return nil
	})
	RegisterDrawApi(app, baseUrl, env)
	RegisterRecordParquetApi(app, baseUrl, env)
}

func RegisterDrawApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewMasterAPI(env, baseUrl))
}

// RegisterRecordParquetApi 注册银行流水 parquet 导入 API
func RegisterRecordParquetApi(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewRecordParquetAPI(env, baseUrl))
}
