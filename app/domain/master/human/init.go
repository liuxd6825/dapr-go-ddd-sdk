package human

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// Init 初始化人员子域，注册 Human 及 13 个子实体的 REST API
func Init(app *iris.Application, baseUrl string, env *env.Env) {
	err := logs.DebugStart(context.Background(), logs.Fields{"service": "human"}, func() error {
		restapi.RegisterHuman(app, baseUrl, env)
		restapi.RegisterHumanAccount(app, baseUrl, env)
		restapi.RegisterHumanAddress(app, baseUrl, env)
		restapi.RegisterHumanCapital(app, baseUrl, env)
		restapi.RegisterHumanCompany(app, baseUrl, env)
		restapi.RegisterHumanContract(app, baseUrl, env)
		restapi.RegisterHumanCredential(app, baseUrl, env)
		restapi.RegisterHumanExt(app, baseUrl, env)
		restapi.RegisterHumanHuman(app, baseUrl, env)
		restapi.RegisterHumanLink(app, baseUrl, env)
		restapi.RegisterHumanProduct(app, baseUrl, env)
		restapi.RegisterHumanRecord(app, baseUrl, env)
		restapi.RegisterHumanReportedAmount(app, baseUrl, env)
		restapi.RegisterHumanSuspectAmount(app, baseUrl, env)
		return nil
	})
	if err != nil {
		panic(err)
	}
}