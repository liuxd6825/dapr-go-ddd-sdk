package company

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/restapi"
	masterconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化公司子域，注册 Company 及 6 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *envs.Env, config masterconfigs.Config) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service-name": "company"}, func() error {
		restapi.RegisterCompany(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyAccount(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyCompany(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyContract(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyHuman(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyProduct(app, baseUrl, env, config.MongoDbKey)
		restapi.RegisterCompanyRecord(app, baseUrl, env, config.MongoDbKey)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
