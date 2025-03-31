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
	return a.TenantId
}
func (a *Aggregate) GetAggId() string {
	return a.AggId
}
func (a *Aggregate) GetAggType() string {
	return a.AggType
}
func (a *Aggregate) GetAggVer() string {
	return a.AggVer
}
