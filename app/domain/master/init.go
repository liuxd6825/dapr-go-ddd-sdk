package master

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

func Init(app *iris.Application, baseUrl string, env *envs.Env) {
	masterconfigs.Init(env)
	record.Init(app, baseUrl, env)
	company.Init(app, baseUrl, env, masterconfigs.GetConfig())
	contract.Init(app, baseUrl, env)
	human.Init(app, baseUrl, env)
	product.Init(app, baseUrl, env)
}
