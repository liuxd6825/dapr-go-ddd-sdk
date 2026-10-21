package graph

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/configs"
	draw_restapi "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/restapi"
	master_restapi "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

func RegisterAllApi(app *iris.Application, baseUrl string, env *envs.Env) {
	graphconfigs.Init(env)
	draw_restapi.RegisterAllApi(app, baseUrl, env)
	master_restapi.RegisterAllApi(app, baseUrl, env)
}
