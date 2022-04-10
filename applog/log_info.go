package applog

type LogInfo struct {
	TenantId  string
	ClassName string
	FuncName  string
	Level     Level
	Message   string
}

func NewLogInfo(tenantId, className, funcName, message string, level Level) *LogInfo {
	return &LogInfo{
		TenantId:  tenantId,
		ClassName: className,
		FuncName:  funcName,
		Message:   message,
		Level:     level,
	}
}

func (i *LogInfo) GetClassName() string {
	return i.ClassName
}
func (i *LogInfo) GetTenantId() string {
	return i.TenantId
}
func (i *LogInfo) GetFuncName() string {
	return i.FuncName
}
func (i *LogInfo) GetLevel() Level {
	return i.Level
}
func (i *LogInfo) GetMessage() string {
	return i.Message
}
