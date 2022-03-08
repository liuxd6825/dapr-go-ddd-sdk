package test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"testing"
)

func TestEventStorage_LoadAggregate(t *testing.T) {
	eventStorage := NewEventStorage()

	agg, _, err := eventStorage.LoadAggregate(context.Background(), "tenant_1", "001", &TestAggregate{})
	if err != nil {
		panic(err)
	}
	println(agg)
}

func TestEventStorage_LoadEvents(t *testing.T) {
	sidecar := NewEventStorage()
	req := &ddd.LoadEventsRequest{
		TenantId:    "tenant_1",
		AggregateId: "001",
	}
	respData, err := sidecar.LoadEvents(context.Background(), req)
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
	req := &ddd.ApplyEventRequest{
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
	respData, err := sidecar.ApplyEvent(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_SaveSnapshot(t *testing.T) {
	sidecar := NewEventStorage()
	req := &ddd.SaveSnapshotRequest{
		TenantId:          "tenantId_1",
		AggregateId:       "aggregateId_001",
		AggregateType:     "system.user",
		Metadata:          map[string]interface{}{"token": "token", "user": "user"},
		AggregateData:     map[string]interface{}{"userId": "001", "userName": "lxd"},
		AggregateRevision: "1.0",
		SequenceNumber:    1,
	}
	respData, err := sidecar.SaveSnapshot(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func NewEventStorage() ddd.EventStorage {
	return ddd.NewEventStorage("localhost", 3500, "pubsub")
}
