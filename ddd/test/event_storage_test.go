package test

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"testing"
	"time"
)

func TestEventStorage_LoadAggregate(t *testing.T) {
	eventStorage, err := NewEventStorage()
	agg, _, err := eventStorage.LoadAggregate(context.Background(), "tenant_1", "001", &TestAggregate{})
	if err != nil {
		panic(err)
	}
	println(agg)
}

func TestEventStorage_LoadEvents(t *testing.T) {
	eventStorage, err := NewEventStorage()
	req := &ddd.LoadEventsRequest{
		TenantId:    "tenant_1",
		AggregateId: "001",
	}
	respData, err := eventStorage.LoadEvents(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_ApplyEvent(t *testing.T) {
	eventStorage, err := NewEventStorage()
	id := newId()
	req := &ddd.ApplyEventRequest{
		TenantId:      "tenantId_1",
		CommandId:     id,
		EventId:       id,
		Metadata:      map[string]string{"token": "token", "user": "user"},
		EventData:     map[string]interface{}{"userId": "001", "userName": "lxd"},
		EventRevision: "1.0",
		EventType:     "CreateUserEvent",
		AggregateId:   "001",
		AggregateType: "system.user",
		PubsubName:    "pubsub",
		Topic:         "topic1",
	}
	respData, err := eventStorage.ApplyEvent(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_SaveSnapshot(t *testing.T) {
	eventStorage, err := NewEventStorage()
	req := &ddd.SaveSnapshotRequest{
		TenantId:          "tenantId_1",
		AggregateId:       "aggregateId_001",
		AggregateType:     "system.user",
		Metadata:          map[string]interface{}{"token": "token", "user": "user"},
		AggregateData:     map[string]interface{}{"userId": "001", "userName": "lxd"},
		AggregateRevision: "1.0",
		SequenceNumber:    1,
	}
	respData, err := eventStorage.SaveSnapshot(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func NewEventStorage() (ddd.EventStorage, error) {
	return ddd.NewDaprEventStorage("localhost", 3500, ddd.PubsubName("pubsub"))
}

func newId() string {
	return fmt.Sprintf("%d", time.Now().Nanosecond())
}
