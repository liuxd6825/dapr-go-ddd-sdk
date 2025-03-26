package userlog

type aggregate struct {
	TenantId         string
	AggregateId      string
	AggregateType    string
	AggregateVersion string
}

const AggregateType = "system.UserLog"
const SystemTenantId = "system"

func newAggregate(tenantId, userId string) *aggregate {
	return &aggregate{
		TenantId:         SystemTenantId,
		AggregateId:      newAggregateId(userId),
		AggregateType:    AggregateType,
		AggregateVersion: "v1",
	}
}

func newAggregateId(userId string) string {
	return userId + "(UserLog)"
}

func (a *aggregate) GetTenantId() string {
	return a.TenantId
}

func (a *aggregate) GetAggregateId() string {
	return a.AggregateId
}

func (a *aggregate) GetAggregateType() string {
	return a.AggregateType
}

func (a *aggregate) GetAggregateVersion() string {
	return a.AggregateVersion
}
