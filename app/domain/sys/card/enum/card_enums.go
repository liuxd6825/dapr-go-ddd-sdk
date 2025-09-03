package enum

type HomeType string

const (
	SystemHome HomeType = "SystemHome" // 系统首页
	UserHome   HomeType = "UserHome"   // 用户首页
)

func (e HomeType) Title() string {
	res := "UNKNOWN"
	switch e {
	case SystemHome:
		res = "系统首页"
	case UserHome:
		res = "用户首页"
	default:
		return "UNKNOWN"
	}
	return res
}

type CardLevel string

const (
	AppLevel  CardLevel = "AppLevel"  // 应用级
	UserLevel CardLevel = "UserLevel" // 用户级
)

func (e CardLevel) Title() string {
	res := "UNKNOWN"
	switch e {
	case AppLevel:
		res = "应用级"
	case UserLevel:
		res = "用户级"
	default:
		return "UNKNOWN"
	}
	return res
}

type CardMode string

const (
	Card CardMode = "Card" // 卡片
	Link CardMode = "Link" // 链接
)

func (e CardMode) Title() string {
	res := "UNKNOWN"
	switch e {
	case Card:
		res = "卡片"
	case Link:
		res = "链接"
	default:
		return "UNKNOWN"
	}
	return res
}

type OpenedTarget string

const (
	BlankOpened OpenedTarget = "BlankOpened" // 空白页面
	SelfOpened  OpenedTarget = "SelfOpened"  // 当前页面
)

func (e OpenedTarget) Title() string {
	res := "UNKNOWN"
	switch e {
	case BlankOpened:
		res = "空白页"
	case SelfOpened:
		res = "当前页"
	default:
		return "UNKNOWN"
	}
	return res
}
