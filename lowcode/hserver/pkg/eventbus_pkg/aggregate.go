package eventbus_pkg

type Aggregate struct {
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
