package ddd

import (
	"fmt"
	"github.com/google/uuid"
	"liuxd/dapr/ddd-iris/demo-command-service/infrastructure/ddd/example"
	"testing"
)

func TestEventStorage_LoadAggregate(t *testing.T) {
	sidecar := NewEventStorage()
	agg, _, err := sidecar.LoadAggregate("tenant_1", "001", &example.TestAggregate{})
	if err != nil {
		panic(err)
	}
	println(agg)
}

func TestEventStorage_LoadEvents(t *testing.T) {
	sidecar := NewEventStorage()
	req := &LoadEventsRequest{
		TenantId:    "tenant_1",
		AggregateId: "001",
	}
	respData, err := sidecar.LoadEvents(req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_ApplyEvent(t *testing.T) {
	sidecar := NewEventStorage()
	id, _ := uuid.NewUUID()
	req := &ApplyEventRequest{
		TenantId:      "tenantId_1",
		CommandId:     id.String(),
		EventId:       id.String(),
		Metadata:      map[string]string{"token": "token", "user": "user"},
		EventData:     map[string]interface{}{"userId": "001", "userName": "lxd"},
		EventRevision: "1.0",
		EventType:     "CreateUserEvent",
		AggregateId:   "001",
		AggregateType: "system.user",
		PubsubName:    "pubsub",
		Topic:         "topic1",
	}
	respData, err := sidecar.ApplyEvent(req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_SaveSnapshot(t *testing.T) {
	sidecar := NewEventStorage()
	req := &SaveSnapshotRequest{
		TenantId:          "tenantId_1",
		AggregateId:       "aggregateId_001",
		AggregateType:     "system.user",
		Metadata:          map[string]interface{}{"token": "token", "user": "user"},
		AggregateData:     map[string]interface{}{"userId": "001", "userName": "lxd"},
		AggregateRevision: "1.0",
		SequenceNumber:    1,
	}
	respData, err := sidecar.SaveSnapshot(req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func NewEventStorage() EventStorage {
	return newEventStorage("localhost", 3500, "pubsub")
}
