package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RegisterContract 注册合同 REST 控制器
func RegisterContract(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractAPI(baseUrl))
}

// RegisterContractCompany 注册合同-公司 REST 控制器
func RegisterContractCompany(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractCompanyAPI(baseUrl))
}

// RegisterContractContract 注册合同-合同 REST 控制器
func RegisterContractContract(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractContractAPI(baseUrl))
}

// RegisterContractHuman 注册合同-人员 REST 控制器
func RegisterContractHuman(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractHumanAPI(baseUrl))
}

// RegisterContractProduct 注册合同-产品 REST 控制器
func RegisterContractProduct(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractProductAPI(baseUrl))
}

// RegisterContractRecord 注册合同-记录 REST 控制器
func RegisterContractRecord(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewContractRecordAPI(baseUrl))
}