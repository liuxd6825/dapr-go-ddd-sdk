package common

type Aggregate struct {
	TenantId         string `json:"tenantId,omitempty"`
	AggregateId      string `json:"aggregateId,omitempty"`
	AggregateType    string `json:"aggregateType,omitempty"`
	AggregateVersion string `json:"aggregateVersion,omitempty"`
}

func NewAggregate() *Aggregate {
	return &Aggregate{}
}

func (a *Aggregate) GetTenantId() string {
	return a.TenantId
}

func (a *Aggregate) GetAggregateId() string {
	return a.AggregateId
}

func (a *Aggregate) GetAggregateType() string {
	return a.AggregateType
}

func (a *Aggregate) GetAggregateVersion() string {
	return a.AggregateVersion
}

func (a *Aggregate) SetTenantId(tenantId string) {
	a.TenantId = tenantId
}

func (a *Aggregate) SetAggregateId(aggregateId string) {
	a.AggregateId = aggregateId
}

func (a *Aggregate) SetAggregateType(aggregateType string) {
	a.AggregateType = aggregateType
}

func (a *Aggregate) SetAggregateVersion(aggregateVersion string) {
	a.AggregateVersion = aggregateVersion
}
