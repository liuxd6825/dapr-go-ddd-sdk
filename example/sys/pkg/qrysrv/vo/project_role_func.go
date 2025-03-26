package vo

// ProjectRoleFunctionVo 对应Java的ProjectRoleFunctionVo类，继承BaseVo
type ProjectRoleFunctionVo struct {
	BaseVo               // 继承BaseVo
	TenantID      string `gorm:"type:varchar(36)" json:"tenantId"`      // 租户ID
	ProjectRoleID string `gorm:"type:varchar(36)" json:"projectRoleId"` // 项目角色ID
	ModuleType    string `gorm:"type:varchar(50)" json:"moduleType"`    // 模块类型
	ModuleID      string `gorm:"type:varchar(36)" json:"moduleId"`      // 模块ID
	FunctionID    string `gorm:"type:varchar(36)" json:"functionId"`    // 功能ID
}
