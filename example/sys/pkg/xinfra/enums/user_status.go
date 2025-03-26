package enums

type UserStatus int

const (
	UserStatus_REVIEW UserStatus = iota
	UserStatus_USEING
	UserStatus_DISABLED
)

func (u UserStatus) String() string {
	switch u {
	case UserStatus_REVIEW:
		return "审核中"
	case UserStatus_USEING:
		return "已启用"
	case UserStatus_DISABLED:
		return "已禁用"
	}
	return "NA"
}
