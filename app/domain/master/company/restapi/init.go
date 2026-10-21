package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RegisterCompany 注册公司 REST 控制器
func RegisterCompany(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyAPI(baseUrl, dbKey))
}

// RegisterCompanyAccount 注册公司账号 REST 控制器
func RegisterCompanyAccount(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyAccountAPI(baseUrl, dbKey))
}

// RegisterCompanyCompany 注册公司关系 REST 控制器
func RegisterCompanyCompany(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyCompanyAPI(baseUrl, dbKey))
}

// RegisterCompanyContract 注册公司合同 REST 控制器
func RegisterCompanyContract(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyContractAPI(baseUrl, dbKey))
}

// RegisterCompanyHuman 注册公司人员 REST 控制器
func RegisterCompanyHuman(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyHumanAPI(baseUrl, dbKey))
}

// RegisterCompanyProduct 注册公司产品 REST 控制器
func RegisterCompanyProduct(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyProductAPI(baseUrl, dbKey))
}

// RegisterCompanyRecord 注册公司流水 REST 控制器（只读）
func RegisterCompanyRecord(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewCompanyRecordAPI(baseUrl, dbKey))
}
