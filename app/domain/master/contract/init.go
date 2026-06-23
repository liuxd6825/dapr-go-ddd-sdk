package contract

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化合同子域，注册 Contract 及 5 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service": "contract"}, func() error {
		restapi.RegisterContract(app, baseUrl, env)
		restapi.RegisterContractCompany(app, baseUrl, env)
		restapi.RegisterContractContract(app, baseUrl, env)
		restapi.RegisterContractHuman(app, baseUrl, env)
		restapi.RegisterContractProduct(app, baseUrl, env)
		restapi.RegisterContractRecord(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}