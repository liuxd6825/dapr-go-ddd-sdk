package vo

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sysmanage/pkg/xinfra/enums"
)

// FunctionVo 对应Java的FunctionVo类，继承BaseVo并适配GORM
type FunctionVo struct {
	BaseVo                         // 继承BaseVo
	ModuleType    enums.ModuleType `gorm:"type:varchar(50)" json:"moduleType"`    // 模块类型
	ModuleID      string           `gorm:"type:varchar(36)" json:"moduleId"`      // 模块ID
	FunctionName  string           `gorm:"type:varchar(255)" json:"functionName"` // 功能名称
	FunctionDesc  string           `gorm:"type:varchar(512)" json:"functionDesc"` // 功能描述
	FunctionOrder int              `gorm:"type:int" json:"functionOrder"`         // 功能顺序
	Status        enums.FuncStatus `gorm:"type:varchar(50)" json:"status"`        // 状态枚举
}
