package analysis

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "analyse"}, func() error {
		analysisconfigs.Init(env)
		restapi.Register(app, baseUrl, env)
		service.RegisterWorkflow()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
