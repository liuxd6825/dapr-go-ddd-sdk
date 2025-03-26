package enums

type FuncStatus int

const (
	FuncStatus_Review    FuncStatus = iota //"审核中"
	FuncStatus_Effective                   //"已生效"
	FuncStatus_Logout                      // "已注销"
)

func (f FuncStatus) String() string {
	switch f {
	case FuncStatus_Review:
		return "审核中"
	case FuncStatus_Effective:
		return "已生效"
	case FuncStatus_Logout:
		return "已注销"
	default:
		return "NA"
	}
}
