package enum

type TenantStatus int

const (
	Review    TenantStatus = 0 // 审核中
	Effective TenantStatus = 1 // 已生效
	Stop      TenantStatus = 2 // 已停用
	Logout    TenantStatus = 3 // 已注销
)

func (e TenantStatus) String() string {
	res := "UNKNOWN"
	switch e {
	case Review:
		res = "Review"
	case Effective:
		res = "Effective"
	case Stop:
		res = "Stop"
	case Logout:
		res = "Logout"
	default:
		res = "UNKNOWN"
	}
	return res
}

func (e TenantStatus) Title() string {
	res := "UNKNOWN"
	switch e {
	case Review:
		res = "审核中"
	case Effective:
		res = "已生效"
	case Stop:
		res = "已停用"
	case Logout:
		res = "已注销"
	default:
		return "UNKNOWN"
	}
	return res
}

type UserStatus int

const (
	Using    UserStatus = 1 // 已启用
	Disabled UserStatus = 0 // 已禁用
)

func (e UserStatus) String() string {
	res := "UNKNOWN"
	switch e {
	case Using:
		res = "Using"
	case Disabled:
		res = "Disabled"
	default:
		res = "UNKNOWN"
	}
	return res
}

func (e UserStatus) Title() string {
	res := "UNKNOWN"
	switch e {
	case Using:
		res = "已启用"
	case Disabled:
		res = "已禁用"
	default:
		return "UNKNOWN"
	}
	return res
}
