package vo

// ProjectRoleVo 对应Java的ProjectRoleVo类，继承BaseVo
type ProjectRoleVo struct {
	BaseVo                         // 继承BaseVo
	TenantID        string         `gorm:"type:varchar(36)" json:"tenantId"` // 租户ID
	Name            string         `gorm:"type:varchar(255)" json:"name"`    // 角色名称
	OrderBy         float64        `gorm:"type:double" json:"orderBy"`       // 排序字段
	ModuleFunctions map[string]Set `gorm:"-" json:"moduleFunctions"`         // 模块功能映射
}

// Set 自定义类型，用于表示Java中的Set
type Set map[string]struct{}
