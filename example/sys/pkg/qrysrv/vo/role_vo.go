package vo

type RoleVo struct {
	BaseVo            // 继承BaseVo
	TenantID  string  `gorm:"type:varchar(36)" json:"tenantId" elastic:"type:keyword"` // 租户ID
	Name      string  `gorm:"type:varchar(255)" json:"name"`                           // 角色名称
	Code      string  `gorm:"type:varchar(255)" json:"code"`                           // 角色代码
	Status    string  `gorm:"type:varchar(50)" json:"status" elastic:"type:keyword"`   // 状态
	RoleOrder float64 `gorm:"type:double" json:"roleOrder"`                            // 角色顺序
}
