package ddd

import "context"

type Aggregate interface {
	OnSourceEvent(ctx context.Context, domainEvent DomainEvent) error
	OnCommand(ctx context.Context, cmd DomainCommand) error
	CreateDomainEvent(ctx context.Context, eventRecord *EventRecord) DomainEvent
	GetAggregateRevision() string
	GetAggregateType() string
	GetAggregateId() string
	GetTenantId() string
	SetTenantId(tenantId string)
}

type AggregateRoot struct {
	Id       string `json:"id"`
	TenantId string `json:"tenantId"`
}

type DomainEventFactory interface {
	NewDomainEvent(eventType string) DomainEvent
}
