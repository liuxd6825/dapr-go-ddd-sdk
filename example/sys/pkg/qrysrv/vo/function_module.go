package vo

// FunctionModuleVo 对应Java的FunctionModuleVo类，继承BaseVo并适配GORM
type FunctionModuleVo struct {
	BaseVo                       // 继承BaseVo
	ModuleType  ModuleType       `gorm:"type:varchar(50)" json:"moduleType"`   // 模块类型
	ModuleName  string           `gorm:"type:varchar(255)" json:"moduleName"`  // 模块名称
	ModuleDesc  string           `gorm:"type:varchar(512)" json:"moduleDesc"`  // 模块描述
	Status      StatusEnum       `gorm:"type:varchar(50)" json:"status"`       // 状态枚举
	ModuleOrder int              `gorm:"type:int" json:"moduleOrder"`          // 模块顺序
	Functions   []FunctionEntity `gorm:"foreignKey:ModuleID" json:"functions"` // 关联的功能实体列表
}
