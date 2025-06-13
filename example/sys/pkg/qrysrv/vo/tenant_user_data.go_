package vo

type TenantUserDataVo struct {
	BaseVo                     // 继承BaseVo
	TenantID      string       `gorm:"type:varchar(36)" json:"tenantId"`       // 租户ID
	TenantName    string       `gorm:"type:varchar(255)" json:"tenantName"`    // 租户名称
	TenantAccount string       `gorm:"type:varchar(255)" json:"tenantAccount"` // 租户账号
	TenantStatus  TenantStatus `gorm:"type:varchar(50)" json:"tenantStatus"`   // 租户状态
	UserID        string       `gorm:"type:varchar(36)" json:"userId"`         // 用户ID
	Account       string       `gorm:"type:varchar(255)" json:"account"`       // 用户账号
	UserStatus    UserStatus   `gorm:"type:varchar(50)" json:"userStatus"`     // 用户状态
	DefaultTenant bool         `gorm:"type:tinyint(1)" json:"defaultTenant"`   // 是否默认租户
	UserType      UserType     `gorm:"type:varchar(50)" json:"userType"`       // 用户类型
}
