package enums

type TenantStatus int

const (
	TenantStatus_REVIEW TenantStatus = iota
	TenantStatus_EFFECTIVE
	TenantStatus_LOGOUT
)

func (t TenantStatus) String() string {
	switch t {
	case TenantStatus_REVIEW:
		return "审核中"
	case TenantStatus_EFFECTIVE:
		return "已生效"
	case TenantStatus_LOGOUT:
		return "已注销"
	}
	return "NA"
}
