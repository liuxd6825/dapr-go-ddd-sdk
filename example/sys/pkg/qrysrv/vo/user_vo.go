package vo

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sysmanage/pkg/xinfra/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type UserVo struct {
	BaseVo                        // 继承BaseVo
	Account      string           `gorm:"type:varchar(255)" json:"account" elastic:"type:keyword"`                     // 账号
	Password     string           `gorm:"type:varchar(255)" json:"password" elastic:"type:text"`                       // 密码
	Name         string           `gorm:"type:varchar(255)" json:"name" elastic:"type:keyword"`                        // 姓名
	Gender       enums.Gender     `gorm:"type:varchar(10)" json:"gender" elastic:"type:keyword"`                       // 性别
	Phone        string           `gorm:"type:varchar(20)" json:"phone" elastic:"type:keyword"`                        // 电话
	Email        string           `gorm:"type:varchar(255)" json:"email" elastic:"type:keyword"`                       // 邮箱
	Work         string           `gorm:"type:varchar(255)" json:"work" elastic:"type:text"`                           // 工作
	Address      string           `gorm:"type:varchar(512)" json:"address" elastic:"type:text"`                        // 地址
	Regdate      *times.Time      `gorm:"type:datetime" json:"regdate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 注册日期
	Status       enums.UserStatus `gorm:"type:varchar(50)" json:"status" elastic:"type:keyword"`                       // 用户状态
	HeadPicture  string           `gorm:"type:varchar(512)" json:"headPicture" elastic:"type:text"`                    // 头像
	ClientID     string           `gorm:"type:varchar(255)" json:"clientId" elastic:"type:text"`                       // 客户端ID
	ClientSecret string           `gorm:"type:varchar(255)" json:"clientSecret" elastic:"type:text"`                   // 客户端密钥
	RedirectUris string           `gorm:"type:varchar(512)" json:"redirectUris" elastic:"type:text"`                   // 重定向URI
	Admin        bool             `gorm:"type:tinyint(1)" json:"admin" elastic:"type:boolean"`                         // 是否管理员
}
