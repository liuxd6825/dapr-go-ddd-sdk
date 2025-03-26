package vo

type UserMenuRouterLinkActiveOptions struct {
	Exact bool `json:"exact"` // 是否精确匹配
}

// UserMenuTreeVo 对应Java的UserMenuTreeVo类
type UserMenuTreeVo struct {
	Type                    string                          `json:"type"`                    // 类型
	Label                   string                          `json:"label"`                   // 标签
	Route                   string                          `json:"route"`                   // 路由
	Icon                    string                          `json:"icon"`                    // 图标
	RouterLinkActiveOptions UserMenuRouterLinkActiveOptions `json:"routerLinkActiveOptions"` // 路由链接激活选项
	Children                []UserMenuTreeVo                `json:"children"`                // 子菜单列表
}
