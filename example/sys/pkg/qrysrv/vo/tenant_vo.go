package vo

// TenantVo 对应Java的TenantVo类，继承BaseVo
type TenantVo struct {
	BaseVo                     // 继承BaseVo
	EffectDate    *types.Time  `gorm:"type:datetime" json:"effectDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 生效日期
	LogoutDate    *types.Time  `gorm:"type:datetime" json:"logoutDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 注销日期
	ExpireDate    *types.Time  `gorm:"type:datetime" json:"expireDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 过期日期
	TenantStatus  TenantStatus `gorm:"type:varchar(50)" json:"tenantStatus" elastic:"type:keyword"`                    // 租户状态
	TenantAddress string       `gorm:"type:varchar(512)" json:"tenantAddress" elastic:"type:text"`                     // 租户地址
	TenantEmail   string       `gorm:"type:varchar(255)" json:"tenantEmail" elastic:"type:text"`                       // 租户邮箱
	TenantName    string       `gorm:"type:varchar(255)" json:"tenantName" elastic:"type:text"`                        // 租户名称
	TenantPhone   string       `gorm:"type:varchar(20)" json:"tenantPhone" elastic:"type:text"`                        // 租户电话
	TenantType    TenantType   `gorm:"type:varchar(50)" json:"tenantType" elastic:"type:keyword"`                      // 租户类型
	TenantAccount string       `gorm:"type:varchar(255)" json:"tenantAccount" elastic:"type:text"`                     // 租户账号
}
