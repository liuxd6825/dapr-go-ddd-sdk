package vo

type TenantMenuDataVo struct {
	BaseVo          // 继承BaseVo
	TenantID string `gorm:"type:varchar(36)" json:"tenantId"` // 租户ID
	MenuID   string `gorm:"type:varchar(36)" json:"menuId"`   // 菜单ID
}
