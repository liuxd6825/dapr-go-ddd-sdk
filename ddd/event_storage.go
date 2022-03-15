package ddd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error)
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error)
	GetPubsubName() string
	GetHost() string
	GetPort() int
}

type EventStorageOption func(EventStorage)

func PubsubName(pubsubName string) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		s.pubsubName = pubsubName
	}
}

func IdleConnTimeout(idleConnTimeout time.Duration) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.IdleConnTimeout = idleConnTimeout
	}
}

func MaxIdleConns(maxIdleConns int) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.MaxIdleConns = maxIdleConns
	}
}

func MaxIdleConnsPerHost(maxIdleConnsPerHost int) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.MaxIdleConnsPerHost = maxIdleConnsPerHost
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
		return errors.New(cmd.GetAggregateId() + " aggregate root already exists.")
	}
	return aggregate.OnCommand(ctx, cmd)
}

func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand, opts ...LoadAggregateOption) error {
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId(), aggregate, opts...)
	if err != nil {
		return err
	}
	if !find {
		return errors.New(cmd.GetAggregateId() + " aggregate root not fond.")
	}
	return aggregate.OnCommand(ctx, cmd)
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
	domainEvent, err := NewDomainEvent(record)
	if err != nil {
		return errors.New(fmt.Sprintf("Method is not found or not exported."))
	}
	return callEventHandler(ctx, handler, record.EventType, record.EventRevision, domainEvent)
}

func callEventHandler(ctx context.Context, handler interface{}, eventType string, eventRevision string, event interface{}) error {
	v := reflect.ValueOf(handler)
	methodName := getMethodName(eventType, eventRevision)
	method := v.MethodByName(methodName)
	if method.IsValid() {
		p1 := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(event),
		}
		resValues := method.Call(p1)
		for _, v := range resValues {
			if err, ok := v.Interface().(error); ok {
				return err
			}
		}
	} else {
		return errors.New(fmt.Sprintf("Method %s is not found or not exported.", methodName))
	}
	return nil
}

func getMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "_", -1)
	return fmt.Sprintf("On%sV%s", name, ver)
}
