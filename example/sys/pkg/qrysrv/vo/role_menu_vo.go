package vo

// RoleMenuVo 对应Java的RoleMenuVo类，继承BaseVo
type RoleMenuVo struct {
	BaseVo          // 继承BaseVo
	Role   *RoleVo  `json:"role"`  // 角色信息
	Menus  []MenuVo `json:"menus"` // 菜单列表
}
