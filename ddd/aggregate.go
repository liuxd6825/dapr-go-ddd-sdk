package ddd

// Aggregate 聚合根接口类
type Aggregate interface {
	GetAggregateRevision() string
	GetAggregateType() string
	GetAggregateId() string
	GetTenantId() string
}
