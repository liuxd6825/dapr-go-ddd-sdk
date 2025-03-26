package enums

type TenantType int

const (
	TenantType_COMPANY TenantType = iota
	TenantType_PERSONAL
)

func (t TenantType) String() string {
	switch t {
	case TenantType_COMPANY:
		return "公司"
	case TenantType_PERSONAL:
		return "个人"
	}
	return "NA"
}
