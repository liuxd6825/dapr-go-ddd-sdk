package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RegisterCompany 注册公司 REST 控制器
func RegisterCompany(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyAPI(baseUrl))
}

// RegisterCompanyAccount 注册公司账号 REST 控制器
func RegisterCompanyAccount(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyAccountAPI(baseUrl))
}

// RegisterCompanyCompany 注册公司关系 REST 控制器
func RegisterCompanyCompany(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyCompanyAPI(baseUrl))
}

// RegisterCompanyContract 注册公司合同 REST 控制器
func RegisterCompanyContract(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyContractAPI(baseUrl))
}

// RegisterCompanyHuman 注册公司人员 REST 控制器
func RegisterCompanyHuman(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyHumanAPI(baseUrl))
}

// RegisterCompanyProduct 注册公司产品 REST 控制器
func RegisterCompanyProduct(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyProductAPI(baseUrl))
}

// RegisterCompanyRecord 注册公司流水 REST 控制器（只读）
func RegisterCompanyRecord(app *iris.Application, baseUrl string, env *env.Env) {
	restapi.RegisterController(app, NewCompanyRecordAPI(baseUrl))
}