package product

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化产品子域，注册 Product 及 5 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service": "product"}, func() error {
		restapi.RegisterProduct(app, baseUrl, env)
		restapi.RegisterProductCompany(app, baseUrl, env)
		restapi.RegisterProductContract(app, baseUrl, env)
		restapi.RegisterProductHuman(app, baseUrl, env)
		restapi.RegisterProductProduct(app, baseUrl, env)
		restapi.RegisterProductRecord(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}