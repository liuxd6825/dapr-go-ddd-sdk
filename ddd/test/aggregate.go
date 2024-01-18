package test

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type TestAggregate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenantId"`
}

func (t TestAggregate) OnSourceEvent(ctx context.Context, domainEvent ddd.DomainEvent) error {
	return nil
}

func (t TestAggregate) OnCommand(ctx context.Context, cmd ddd.Command) error {
	return nil
}

func (t TestAggregate) CreateDomainEvent(ctx context.Context, eventRecord *dapr.EventRecord) ddd.DomainEvent {
	return nil
}

func (t TestAggregate) GetAggregateVersion() string {
	return "1.0"
}

func (t TestAggregate) GetAggregateType() string {
	return "TestAggregateType"
}

func (t TestAggregate) GetAggregateId() string {
	return t.Id
}

func (t TestAggregate) GetTenantId() string {
	return t.TenantId
}

func (t TestAggregate) SetTenantId(tenantId string) {
	t.TenantId = tenantId
}
