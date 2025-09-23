package analysis

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func Init(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "analyse"}, func() error {
		restapi.Register(app, baseUrl, env)
		//service.RegisterWorkflow()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
