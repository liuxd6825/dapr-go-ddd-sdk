package ddd

import (
	"context"
	"errors"
	"strings"
)

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error)
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error)
}

func NewEventStorage(host string, port int, options ...func(s EventStorage)) (EventStorage, error) {
	return NewDaprEventStorage(host, port, options...)
}

func newEventStorage(host string, port int, options ...func(s EventStorage)) (EventStorage, error) {
	return NewDaprEventStorage(host, port, options...)
}

func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error) {
	return _eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
}

func LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	return _eventStorage.LoadEvents(ctx, req)
}

func Apply(ctx context.Context, pubsubName string, aggregate Aggregate, event DomainEvent, metadata map[string]string) error {
	topic := strings.ToLower(event.GetEventType())
	req := &ApplyEventRequest{
		TenantId:      event.GetTenantId(),
		CommandId:     event.GetCommandId(),
		EventId:       event.GetEventId(),
		EventRevision: event.GetEventRevision(),
		EventType:     event.GetEventType(),
		AggregateId:   event.GetAggregateId(),
		AggregateType: aggregate.GetAggregateType(),
		Metadata:      metadata,
		EventData:     event,
		PubsubName:    pubsubName,
		Topic:         topic,
	}
	if _, err := _eventStorage.ApplyEvent(ctx, req); err != nil {
		return err
	}
	if err := aggregate.OnSourceEvent(ctx, event); err != nil {
		return err
	}
	return nil
}

func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand) error {
	ok, err := _eventStorage.ExistAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId())
	if err != nil {
		return err
	}
	if ok {
		return errors.New(cmd.GetAggregateId() + " aggregate root already exists.")
	}
	return aggregate.OnCommand(ctx, cmd)
}

func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand) error {
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId(), aggregate)
	if err != nil {
		return err
	}
	if !find {
		return errors.New(cmd.GetAggregateId() + " aggregate root not fond.")
	}
	return aggregate.OnCommand(ctx, cmd)
}

func applyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error) {
	return _eventStorage.ApplyEvent(ctx, req)
}

func saveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	return _eventStorage.SaveSnapshot(ctx, req)
}
