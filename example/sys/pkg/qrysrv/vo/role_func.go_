package vo

// RoleFunctionVo 对应Java的RoleFunctionVo类，继承BaseVo
type RoleFunctionVo struct {
	BaseVo            // 继承BaseVo
	TenantID   string `gorm:"type:varchar(36)" json:"tenantId" elastic:"type:keyword"`   // 租户ID
	RoleID     string `gorm:"type:varchar(36)" json:"roleId" elastic:"type:keyword"`     // 角色ID
	ModuleType string `gorm:"type:varchar(50)" json:"moduleType" elastic:"type:keyword"` // 模块类型
	ModuleID   string `gorm:"type:varchar(36)" json:"moduleId" elastic:"type:keyword"`   // 模块ID
	FunctionID string `gorm:"type:varchar(36)" json:"functionId" elastic:"type:keyword"` // 功能ID
}
