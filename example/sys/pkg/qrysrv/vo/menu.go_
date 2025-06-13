package vo

// MenuVo 对应Java的MenuVo类，继承BaseVo并适配GORM和ES
type MenuVo struct {
	BaseVo                               // 继承BaseVo
	ParentID                     string  `gorm:"type:varchar(36)" json:"parentId" elastic:"type:keyword"`                   // 父级ID
	Name                         string  `gorm:"type:varchar(255)" json:"name" elastic:"type:text"`                         // 名称
	Code                         string  `gorm:"type:varchar(255)" json:"code" elastic:"type:text"`                         // 代码
	Route                        string  `gorm:"type:varchar(255)" json:"route" elastic:"type:text"`                        // 路由
	MenuType                     string  `gorm:"type:varchar(50)" json:"menuType" elastic:"type:text"`                      // 菜单类型
	Icon                         string  `gorm:"type:varchar(255)" json:"icon" elastic:"type:text"`                         // 图标
	SubApplication               string  `gorm:"type:varchar(255)" json:"subApplication" elastic:"type:text"`               // 子应用
	SubApplicationName           string  `gorm:"type:varchar(255)" json:"subApplicationName" elastic:"type:text"`           // 子应用名称
	SubApplicationPrefix         string  `gorm:"type:varchar(255)" json:"subApplicationPrefix" elastic:"type:text"`         // 子应用前缀
	Status                       string  `gorm:"type:varchar(50)" json:"status" elastic:"type:keyword"`                     // 状态
	MenuOrder                    float64 `gorm:"type:double" json:"menuOrder" elastic:"type:double"`                        // 菜单顺序
	RouterLinkActiveOptionsExact string  `gorm:"type:varchar(255)" json:"routerLinkActiveOptionsExact" elastic:"type:text"` // 路由链接激活选项
	Container                    string  `gorm:"type:varchar(255)" json:"container" elastic:"type:text"`                    // 容器
	ActiveRule                   string  `gorm:"type:varchar(255)" json:"activeRule" elastic:"type:text"`                   // 激活规则
}
