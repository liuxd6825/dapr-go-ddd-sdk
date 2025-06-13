package vo

type TenantUserVo struct {
	BaseVo                     // 继承BaseVo
	TenantID      string       `gorm:"type:varchar(36)" json:"tenantId" elastic:"type:keyword"`                     // 租户ID
	TenantName    string       `gorm:"type:varchar(255)" json:"tenantName" elastic:"type:text"`                     // 租户名称
	TenantAccount string       `gorm:"type:varchar(255)" json:"tenantAccount" elastic:"type:text"`                  // 租户账号
	TenantStatus  TenantStatus `gorm:"type:varchar(50)" json:"tenantStatus" elastic:"type:keyword"`                 // 租户状态
	DefaultTenant bool         `gorm:"type:tinyint(1)" json:"defaultTenant" elastic:"type:boolean"`                 // 是否默认租户
	UserType      UserType     `gorm:"type:varchar(50)" json:"userType" elastic:"type:keyword"`                     // 用户类型
	UserID        string       `gorm:"type:varchar(36)" json:"userId" elastic:"type:keyword"`                       // 用户ID
	Account       string       `gorm:"type:varchar(255)" json:"account" elastic:"type:text"`                        // 用户账号
	Password      string       `gorm:"type:varchar(255)" json:"password" elastic:"type:text"`                       // 密码
	Name          string       `gorm:"type:varchar(255)" json:"name" elastic:"type:text"`                           // 姓名
	Gender        Gender       `gorm:"type:varchar(10)" json:"gender" elastic:"type:keyword"`                       // 性别
	Phone         string       `gorm:"type:varchar(20)" json:"phone" elastic:"type:text"`                           // 电话
	Email         string       `gorm:"type:varchar(255)" json:"email" elastic:"type:text"`                          // 邮箱
	Work          string       `gorm:"type:varchar(255)" json:"work" elastic:"type:text"`                           // 工作
	Address       string       `gorm:"type:varchar(512)" json:"address" elastic:"type:text"`                        // 地址
	Regdate       *types.Time  `gorm:"type:datetime" json:"regdate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 注册日期
	UserStatus    UserStatus   `gorm:"type:varchar(50)" json:"userStatus" elastic:"type:keyword"`                   // 用户状态
	HeadPicture   string       `gorm:"type:varchar(512)" json:"headPicture" elastic:"type:text"`                    // 头像
	Admin         bool         `gorm:"type:tinyint(1)" json:"admin" elastic:"type:boolean"`                         // 是否管理员
}
