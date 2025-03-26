package vo

// LoginUser 对应Java的LoginUser类，适配GORM
type LoginUser struct {
	ID            string     `gorm:"primaryKey;type:varchar(36)" json:"id"`  // 主键
	Account       string     `gorm:"type:varchar(255)" json:"account"`       // 账号
	Regdate       DateTime   `gorm:"type:datetime" json:"regdate"`           // 注册日期
	UserStatus    UserStatus `gorm:"type:varchar(50)" json:"userStatus"`     // 用户状态
	Name          string     `gorm:"type:varchar(255)" json:"name"`          // 姓名
	Phone         string     `gorm:"type:varchar(20)" json:"phone"`          // 电话
	Address       string     `gorm:"type:varchar(512)" json:"address"`       // 地址
	Work          string     `gorm:"type:varchar(255)" json:"work"`          // 工作
	Gender        Gender     `gorm:"type:varchar(10)" json:"gender"`         // 性别
	Email         string     `gorm:"type:varchar(255)" json:"email"`         // 邮箱
	Admin         bool       `gorm:"type:tinyint(1)" json:"admin"`           // 是否管理员
	HeadPicture   string     `gorm:"type:varchar(512)" json:"headPicture"`   // 头像
	DefaultTenant bool       `gorm:"type:tinyint(1)" json:"defaultTenant"`   // 是否默认租户
	UserType      UserType   `gorm:"type:varchar(50)" json:"userType"`       // 用户类型
	TenantID      string     `gorm:"type:varchar(36)" json:"tenantId"`       // 租户ID
	TenantName    string     `gorm:"type:varchar(255)" json:"tenantName"`    // 租户名称
	TenantAccount string     `gorm:"type:varchar(255)" json:"tenantAccount"` // 租户账号
}
