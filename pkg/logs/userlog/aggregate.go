package userlog

type aggregate struct {
	TenantId string
	AggId    string
	AggType  string
	AggVer   string
}

const AggregateType = "system.UserLog"
const SystemTenantId = "system"

func newAggregate(tenantId, userId string) *aggregate {
	return &aggregate{
		TenantId: SystemTenantId,
		AggId:    newAggregateId(userId),
		AggType:  AggregateType,
		AggVer:   "v1",
	}
}

func newAggregateId(userId string) string {
	return userId + "(UserLog)"
}

func (a *aggregate) GetTenantId() string {
	return a.TenantId
}

func (a *aggregate) GetAggId() string {
	return a.AggId
}

func (a *aggregate) GetAggType() string {
	return a.AggType
}

func (a *aggregate) GetAggVer() string {
	return a.AggVer
}
