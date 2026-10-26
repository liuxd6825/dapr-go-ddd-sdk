package product

import (
	"context"

	"github.com/kataras/iris/v12"
	masterconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化产品子域，注册 Product 及 5 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name": "product"}, func() error {
		masterconfigs.Init(env)
		dbKey := masterconfigs.GetMongoDbKey()
		restapi.RegisterProduct(app, baseUrl, env, dbKey)
		restapi.RegisterProductCompany(app, baseUrl, env, dbKey)
		restapi.RegisterProductContract(app, baseUrl, env, dbKey)
		restapi.RegisterProductHuman(app, baseUrl, env, dbKey)
		restapi.RegisterProductProduct(app, baseUrl, env, dbKey)
		restapi.RegisterProductRecord(app, baseUrl, env, dbKey)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
