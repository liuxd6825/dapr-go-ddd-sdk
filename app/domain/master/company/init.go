package company

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化公司子域，注册 Company 及 6 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service": "company"}, func() error {
		restapi.RegisterCompany(app, baseUrl, env)
		restapi.RegisterCompanyAccount(app, baseUrl, env)
		restapi.RegisterCompanyCompany(app, baseUrl, env)
		restapi.RegisterCompanyContract(app, baseUrl, env)
		restapi.RegisterCompanyHuman(app, baseUrl, env)
		restapi.RegisterCompanyProduct(app, baseUrl, env)
		restapi.RegisterCompanyRecord(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}