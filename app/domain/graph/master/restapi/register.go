package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *env.Env) {
	_ = logs.DebugStart(context.Background(), logs.Fields{"service name ": "graph"}, func() error {
		return nil
	})
}
