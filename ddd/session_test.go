package ddd

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

type TestAggregate struct {
	id               string
	tenantId         string
	aggregateType    string
	aggregateVersion string
}

type TestEvent struct {
	commandId    string
	id           string
	tenantId     string
	eventType    string
	eventVersion string
	aggregateId  string
	createdTime  time.Time
	data         interface{}
}

func TestNewSession(t *testing.T) {
	ctx := context.Background()
	sCtx, _ := newContext(ctx, "test")
	val, ok := getSession(sCtx)
	if ok {
		t.Log(val)
	} else {
		t.Error(errors.New("session is error"))
	}
}

func TestStartSession(t *testing.T) {
	ctx := context.Background()
	err := StartSession(ctx, "test", func(sCtx context.Context, session Session) error {
		agg := NewTestAggregate()
		if _, err := CreateEvent(sCtx, agg, NewCreateEvent(agg.id)); err != nil {
			return err
		}
		if _, err := ApplyEvent(sCtx, agg, NewUpdateEvent(agg.id)); err != nil {
			return err
		}
		if _, err := DeleteEvent(sCtx, agg, NewDeleteEvent(agg.id)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func NewTestAggregate() *TestAggregate {
	return &TestAggregate{
		id:               uuid.New().String(),
		tenantId:         "test",
		aggregateType:    "ddd.test",
		aggregateVersion: "1.0",
	}
}

func NewCreateEvent(aggregateId string) *TestEvent {
	id := uuid.NewString()
	event := &TestEvent{
		commandId:    id,
		id:           id,
		aggregateId:  aggregateId,
		tenantId:     "test",
		eventType:    "create",
		eventVersion: "v1.0",
		createdTime:  time.Now(),
	}
	return event
}

func NewUpdateEvent(aggregateId string) *TestEvent {
	id := uuid.NewString()
	event := &TestEvent{
		commandId:    id,
		id:           id,
		aggregateId:  aggregateId,
		tenantId:     "test",
		eventType:    "update",
		eventVersion: "v1.0",
		createdTime:  time.Now(),
	}
	return event
}

func NewDeleteEvent(aggregateId string) *TestEvent {
	id := uuid.NewString()
	event := &TestEvent{
		commandId:    id,
		id:           id,
		aggregateId:  aggregateId,
		tenantId:     "test",
		eventType:    "delete",
		eventVersion: "v1.0",
		createdTime:  time.Now(),
	}
	return event
}

func (a *TestAggregate) GetTenantId() string {
	return a.tenantId
}

func (a *TestAggregate) GetAggregateId() string {
	return a.id
}

func (a *TestAggregate) GetAggregateType() string {
	return a.aggregateType
}

func (a *TestAggregate) GetAggregateVersion() string {
	return a.aggregateVersion
}

func (t *TestEvent) GetTenantId() string {
	return t.tenantId
}

func (t *TestEvent) GetCommandId() string {
	return t.commandId
}

func (t *TestEvent) GetEventId() string {
	return t.id
}

func (t *TestEvent) GetEventType() string {
	return t.eventType
}

func (t *TestEvent) GetEventVersion() string {
	return t.eventVersion
}

func (t *TestEvent) GetAggregateId() string {
	return t.aggregateId
}

func (t *TestEvent) GetCreatedTime() time.Time {
	return t.createdTime
}

func (t *TestEvent) GetData() interface{} {
	return t.data
}
