package vo

// MenuTreeVo 对应Java的MenuTreeVo类，继承BaseVo并适配GORM
type MenuTreeVo struct {
	BaseVo                                    // 继承BaseVo
	ParentID                     string       `gorm:"type:varchar(36)" json:"parentId"`                      // 父级ID
	Name                         string       `gorm:"type:varchar(255)" json:"name"`                         // 名称
	Code                         string       `gorm:"type:varchar(255)" json:"code"`                         // 代码
	Route                        string       `gorm:"type:varchar(255)" json:"route"`                        // 路由
	MenuType                     string       `gorm:"type:varchar(50)" json:"menuType"`                      // 菜单类型
	Icon                         string       `gorm:"type:varchar(255)" json:"icon"`                         // 图标
	SubApplication               string       `gorm:"type:varchar(255)" json:"subApplication"`               // 子应用
	SubApplicationName           string       `gorm:"type:varchar(255)" json:"subApplicationName"`           // 子应用名称
	SubApplicationPrefix         string       `gorm:"type:varchar(255)" json:"subApplicationPrefix"`         // 子应用前缀
	Status                       string       `gorm:"type:varchar(50)" json:"status"`                        // 状态
	MenuOrder                    float64      `gorm:"type:double" json:"menuOrder"`                          // 菜单顺序
	RouterLinkActiveOptionsExact string       `gorm:"type:varchar(255)" json:"routerLinkActiveOptionsExact"` // 路由链接激活选项
	Container                    string       `gorm:"type:varchar(255)" json:"container"`                    // 容器
	ActiveRule                   string       `gorm:"type:varchar(255)" json:"activeRule"`                   // 激活规则
	Children                     []MenuTreeVo `gorm:"-" json:"children"`                                     // 子菜单列表
	HasChildren                  bool         `gorm:"-" json:"hasChildren"`                                  // 是否有子菜单
}
