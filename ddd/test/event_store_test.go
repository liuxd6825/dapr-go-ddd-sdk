package test

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
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
	req := &dapr.LoadEventsRequest{
		TenantId:    "tenant_1",
		AggregateId: "001",
	}
	respData, err := eventStorage.LoadEvent(context.Background(), req)
	if err != nil {
		panic(err)
	}

	if respData != nil {
		fmt.Println(respData)
	}
}

func TestEventStorage_ApplyEvent(t *testing.T) {
	eventStorage, err := NewEventStorage()
	// id := newId()
	req := &dapr.ApplyEventRequest{
		TenantId:      "tenantId_1",
		AggregateId:   "001",
		AggregateType: "system.user",
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
	req := &dapr.SaveSnapshotRequest{
		TenantId:         "tenantId_1",
		AggregateId:      "aggregateId_001",
		AggregateType:    "system.user",
		Metadata:         map[string]string{"token": "token", "user": "user"},
		AggregateData:    map[string]interface{}{"userId": "001", "userName": "lxd"},
		AggregateVersion: "1.0",
		SequenceNumber:   1,
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
	daprDddClient, err := dapr.NewDaprDddClient("localhost", 3500, 0000)
	if err != nil {
		return nil, err
	}
	return ddd.NewGrpcEventStore(daprDddClient, ddd.PubsubName("pubsub"))
}

func newId() string {
	return fmt.Sprintf("%d", time.Now().Nanosecond())
}
