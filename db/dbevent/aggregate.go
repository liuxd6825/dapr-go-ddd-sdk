package dbevent

type Aggregate struct {
	TenantId string `json:"tenantId"`
	AggId    string `json:"aggId"`
	AggType  string `json:"aggType"`
	AggVer   string `json:"aggVer"`
}

func NewAggregate() *Aggregate {
	return &Aggregate{}
}
func (a *Aggregate) GetTenantId() string {
	return ""
}
func (a *Aggregate) GetAggregateId() string {
	return ""
}
func (a *Aggregate) GetAggregateType() string {
	return ""
}
func (a *Aggregate) GetAggregateVersion() string {
	return ""
}
