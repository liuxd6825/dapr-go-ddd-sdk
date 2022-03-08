package ddd

type DomainEvent interface {
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
	GetEventType() string
	GetEventRevision() string
	GetAggregateId() string
}
