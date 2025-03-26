package enums

type UserType int

const (
	UserType_TenantUser UserType = iota
	UserType_TenantAdmin
)

func (u UserType) String() string {
	switch u {
	case UserType_TenantUser:
		return "租户用户"
	case UserType_TenantAdmin:
		return "租户管理员"
	}
	return "NA"
}
