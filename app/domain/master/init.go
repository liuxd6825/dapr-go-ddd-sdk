package master

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func Init(app *iris.Application, baseUrl string, env *env.Env) {
	record.Init(app, baseUrl, env)
	company.Init(app, baseUrl, env)
	contract.Init(app, baseUrl, env)
	human.Init(app, baseUrl, env)
	product.Init(app, baseUrl, env)
}
