package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RegisterProduct 注册产品 REST 控制器
func RegisterProduct(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductAPI(baseUrl))
}

// RegisterProductCompany 注册产品-公司 REST 控制器
func RegisterProductCompany(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductCompanyAPI(baseUrl))
}

// RegisterProductContract 注册产品-合同 REST 控制器
func RegisterProductContract(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductContractAPI(baseUrl))
}

// RegisterProductHuman 注册产品-人员 REST 控制器
func RegisterProductHuman(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductHumanAPI(baseUrl))
}

// RegisterProductProduct 注册产品-产品 REST 控制器
func RegisterProductProduct(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductProductAPI(baseUrl))
}

// RegisterProductRecord 注册产品-记录 REST 控制器
func RegisterProductRecord(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewProductRecordAPI(baseUrl))
}