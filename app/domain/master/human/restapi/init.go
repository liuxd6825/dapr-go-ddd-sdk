package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RegisterHuman 注册人员 REST 控制器
func RegisterHuman(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanAPI(baseUrl, dbKey))
}

// RegisterHumanAccount 注册人员-账户 REST 控制器
func RegisterHumanAccount(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanAccountAPI(baseUrl, dbKey))
}

// RegisterHumanAddress 注册人员-地址 REST 控制器
func RegisterHumanAddress(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanAddressAPI(baseUrl, dbKey))
}

// RegisterHumanCapital 注册人员-资产 REST 控制器
func RegisterHumanCapital(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanCapitalAPI(baseUrl, dbKey))
}

// RegisterHumanCompany 注册人员-公司 REST 控制器
func RegisterHumanCompany(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanCompanyAPI(baseUrl, dbKey))
}

// RegisterHumanContract 注册人员-合同 REST 控制器
func RegisterHumanContract(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanContractAPI(baseUrl, dbKey))
}

// RegisterHumanCredential 注册人员-证件 REST 控制器
func RegisterHumanCredential(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanCredentialAPI(baseUrl, dbKey))
}

// RegisterHumanExt 注册人员-扩展信息 REST 控制器
func RegisterHumanExt(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanExtAPI(baseUrl, dbKey))
}

// RegisterHumanHuman 注册人员-人员 REST 控制器
func RegisterHumanHuman(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanHumanAPI(baseUrl))
}

// RegisterHumanLink 注册人员-联系方式 REST 控制器
func RegisterHumanLink(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanLinkAPI(baseUrl))
}

// RegisterHumanProduct 注册人员-产品 REST 控制器
func RegisterHumanProduct(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanProductAPI(baseUrl, dbKey))
}

// RegisterHumanRecord 注册人员-记录 REST 控制器
func RegisterHumanRecord(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanRecordAPI(baseUrl, dbKey))
}

// RegisterHumanReportedAmount 注册报案人金额 REST 控制器
func RegisterHumanReportedAmount(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanReportedAmountAPI(baseUrl, dbKey))
}

// RegisterHumanSuspectAmount 注册嫌疑人金额 REST 控制器
func RegisterHumanSuspectAmount(app *iris.Application, baseUrl string, env *envs.Env, dbKey string) {
	restapi.RegisterController(app, NewHumanSuspectAmountAPI(baseUrl, dbKey))
}
