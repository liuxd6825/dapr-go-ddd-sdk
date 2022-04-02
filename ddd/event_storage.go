package ddd

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"reflect"
	"strings"
)

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error)
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error)
	GetPubsubName() string
}

type EventStorageOption func(EventStorage)

func PubsubName(pubsubName string) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		s.pubsubName = pubsubName
	}
}

type LoadAggregateOptions struct {
	eventStorageKey string
}
type LoadAggregateOption func(*LoadAggregateOptions)

func LoadAggregateKey(eventStorageKey string) LoadAggregateOption {
	return func(options *LoadAggregateOptions) {
		options.eventStorageKey = eventStorageKey
	}
}

func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate, opts ...LoadAggregateOption) (Aggregate, bool, error) {
	options := &LoadAggregateOptions{
		eventStorageKey: "",
	}
	for _, item := range opts {
		item(options)
	}
	eventStorage, err := GetEventStorage(options.eventStorageKey)
	if err != nil {
		return nil, false, err
	}
	return eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
}

func LoadEvents(ctx context.Context, req *LoadEventsRequest, eventStorageKey string) (*LoadEventsResponse, error) {
	eventStorage, err := GetEventStorage(eventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.LoadEvents(ctx, req)
}

type ApplyOptions struct {
	pubsubName      string
	metadata        map[string]string
	eventStorageKey string
}
type ApplyOption func(*ApplyOptions)

func ApplyPubsubName(pubsubName string) ApplyOption {
	return func(o *ApplyOptions) {
		o.pubsubName = pubsubName
	}
}

func ApplyEventStorageKey(eventStorageKey string) ApplyOption {
	return func(o *ApplyOptions) {
		o.eventStorageKey = eventStorageKey
	}
}

func ApplyMetadata(metadata map[string]string) ApplyOption {
	return func(o *ApplyOptions) {
		o.metadata = metadata
	}
}

func Apply(ctx context.Context, aggregate Aggregate, event DomainEvent, options ...ApplyOption) error {
	appOptions := &ApplyOptions{
		pubsubName:      "",
		metadata:        map[string]string{},
		eventStorageKey: "",
	}
	for _, option := range options {
		option(appOptions)
	}
	eventStorage, err := GetEventStorage(appOptions.eventStorageKey)
	if err != nil {
		return err
	}
	req := &ApplyEventRequest{
		TenantId:      event.GetTenantId(),
		CommandId:     event.GetCommandId(),
		EventId:       event.GetEventId(),
		EventRevision: event.GetEventRevision(),
		EventType:     event.GetEventType(),
		AggregateId:   event.GetAggregateId(),
		AggregateType: aggregate.GetAggregateType(),
		Metadata:      appOptions.metadata,
		PubsubName:    appOptions.pubsubName,
		EventData:     event,
		Topic:         event.GetEventType(),
	}
	if _, err := eventStorage.ApplyEvent(ctx, req); err != nil {
		return err
	}
	if err := callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventRevision(), event); err != nil {
		return err
	}
	return nil
}

type CreateAggregateOptions struct {
	eventStorageKey string
}
type CreateAggregateOption func(*CreateAggregateOptions)

func CreateAggregateKey(eventStorageKey string) CreateAggregateOption {
	return func(options *CreateAggregateOptions) {
		options.eventStorageKey = eventStorageKey
	}
}

func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand, opts ...CreateAggregateOption) error {
	options := &CreateAggregateOptions{
		eventStorageKey: "",
	}
	for _, item := range opts {
		item(options)
	}

	eventStorage, err := GetEventStorage(options.eventStorageKey)
	if err != nil {
		return err
	}
	ok, err := eventStorage.ExistAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId())
	if err != nil {
		return err
	}
	if ok {
		return ddd_errors.NewNotFondAggregateIdError(cmd.GetAggregateId())
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

func callCommandHandler(ctx context.Context, aggregate Aggregate, cmd DomainCommand) error {
	cmdTypeName := reflect.ValueOf(cmd).Elem().Type().Name()
	methodName := fmt.Sprintf("%s", cmdTypeName)
	metadata := ddd_context.GetMetadataContext(ctx)
	return CallMethod(aggregate, methodName, ctx, cmd, metadata)
}
func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand, opts ...LoadAggregateOption) error {
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId(), aggregate, opts...)
	if err != nil {
		return err
	}
	if !find {
		return ddd_errors.NewNotFondAggregateIdError(aggregate.GetAggregateId())
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

func applyEvent(ctx context.Context, req *ApplyEventRequest, eventStorageKey string) (*ApplyEventsResponse, error) {
	eventStorage, err := GetEventStorage(eventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.ApplyEvent(ctx, req)
}

func saveSnapshot(ctx context.Context, req *SaveSnapshotRequest, eventStorageKey string) (*SaveSnapshotResponse, error) {
	eventStorage, err := GetEventStorage(eventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.SaveSnapshot(ctx, req)
}

func CallEventHandler(ctx context.Context, handler interface{}, record *EventRecord) error {
	event, err := NewDomainEvent(record)
	if err != nil {
		return errors.New(fmt.Sprintf("Method is not found or not exported."))
	}
	return callEventHandler(ctx, handler, record.EventType, record.EventRevision, event)
}

func callEventHandler(ctx context.Context, handler interface{}, eventType string, eventRevision string, event interface{}) error {
	methodName := getEventMethodName(eventType, eventRevision)
	return CallMethod(handler, methodName, ctx, event)
}

func getEventMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "s", -1)
	return fmt.Sprintf("On%sV%s", name, ver)
}
