package example

import (
	"liuxd/dapr/ddd-iris/demo-command-service/infrastructure/ddd"
)

type TestAggregate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenantId"`
}

func (t TestAggregate) OnSourceEvent(domainEvent ddd.DomainEvent) error {
	return nil
}

func (t TestAggregate) OnCommand(cmd ddd.DomainCommand) error {
	return nil
}

func (t TestAggregate) CreateDomainEvent(eventRecord *ddd.EventRecord) ddd.DomainEvent {
	return nil
}

func (t TestAggregate) GetAggregateRevision() string {
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
