package enums

type Gender int

const (
	Gender_None   Gender = iota //"未知"
	Gender_Male                 //"男"
	Gender_Female               // "女"
)

func (g Gender) String() string {
	switch g {
	case Gender_None:
		return "审核中"
	case Gender_Male:
		return "已生效"
	case Gender_Female:
		return "已注销"
	default:
		return "NA"
	}
}
