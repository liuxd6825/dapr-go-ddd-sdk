package applog

//
// LogInfo
// @Description: 日志信息
//
type LogInfo struct {
	TenantId  string
	ClassName string
	FuncName  string
	Level     Level
	Message   string
}

//
// NewLogInfo
// @Description: 新建日志信息
// @param tenantId
// @param className
// @param funcName
// @param message
// @param level
// @return *LogInfo
//
func NewLogInfo(tenantId, className, funcName, message string, level Level) *LogInfo {
	return &LogInfo{
		TenantId:  tenantId,
		ClassName: className,
		FuncName:  funcName,
		Message:   message,
		Level:     level,
	}
}

//
// GetClassName
// @Description: 获取结构名称
// @receiver i
// @return string
//
func (i *LogInfo) GetClassName() string {
	return i.ClassName
}

//
// GetTenantId
// @Description: 获取租户名称
// @receiver i
// @return string
//
func (i *LogInfo) GetTenantId() string {
	return i.TenantId
}

//
// GetFuncName
// @Description: 获取方法名称
// @receiver i
// @return string
//
func (i *LogInfo) GetFuncName() string {
	return i.FuncName
}

//
// GetLevel
// @Description: 获取日志级别
// @receiver i
// @return Level
//
func (i *LogInfo) GetLevel() Level {
	return i.Level
}

//
// GetMessage
// @Description: 获取日志内容
// @receiver i
// @return string
//
func (i *LogInfo) GetMessage() string {
	return i.Message
}
