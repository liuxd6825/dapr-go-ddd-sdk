package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
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

func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate, opts ...LoadAggregateOption) (agg Aggregate, isFound bool, err error) {

	logInfo := &applog.LogInfo{
		TenantId:  tenantId,
		ClassName: "ddd",
		FuncName:  "LoadAggregate",
		Message:   fmt.Sprintf("aggregateId=%s", aggregateId),
		Level:     applog.INFO,
	}

	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		options := &LoadAggregateOptions{
			eventStorageKey: "",
		}
		for _, item := range opts {
			item(options)
		}
		eventStorage, e := GetEventStorage(options.eventStorageKey)
		if e != nil {
			agg, isFound, err = nil, false, e
			return agg, err
		}
		agg, isFound, err = eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
		return agg, err
	})
	return
}

func LoadEvents(ctx context.Context, req *LoadEventsRequest, eventStorageKey string) (resp *LoadEventsResponse, err error) {
	logInfo := &applog.LogInfo{
		TenantId:  req.TenantId,
		ClassName: "ddd",
		FuncName:  "LoadAggregate",
		Message:   fmt.Sprintf("%v", req),
		Level:     applog.INFO,
	}
	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		eventStorage, e := GetEventStorage(eventStorageKey)
		if e != nil {
			resp, err = nil, e
			return resp, err
		}
		resp, err = eventStorage.LoadEvents(ctx, req)
		return resp, err
	})
	return
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

func Apply(ctx context.Context, aggregate Aggregate, event DomainEvent, options ...ApplyOption) (err error) {

	logInfo := &applog.LogInfo{
		TenantId:  aggregate.GetTenantId(),
		ClassName: "ddd",
		FuncName:  "LoadAggregate",
		Message:   fmt.Sprintf("%v", aggregate),
		Level:     applog.INFO,
	}
	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		appOptions := &ApplyOptions{
			pubsubName:      "",
			metadata:        map[string]string{},
			eventStorageKey: "",
		}
		for _, option := range options {
			option(appOptions)
		}
		eventStorage, e := GetEventStorage(appOptions.eventStorageKey)
		if e != nil {
			err = e
			return nil, err
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
		if _, err = eventStorage.ApplyEvent(ctx, req); err != nil {
			return nil, err
		}
		if err = callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventRevision(), event); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return
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

func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...CreateAggregateOption) error {
	options := &CreateAggregateOptions{
		eventStorageKey: "",
	}
	for _, item := range opts {
		item(options)
	}
	rootId := cmd.GetAggregateId().RootId()
	eventStorage, err := GetEventStorage(options.eventStorageKey)
	if err != nil {
		return err
	}
	ok, err := eventStorage.ExistAggregate(ctx, cmd.GetTenantId(), rootId)
	if err != nil {
		return err
	}
	if ok {
		return ddd_errors.NewAggregateIdExistsError(rootId)
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

func callCommandHandler(ctx context.Context, aggregate Aggregate, cmd Command) error {
	cmdTypeName := reflect.ValueOf(cmd).Elem().Type().Name()
	methodName := fmt.Sprintf("%s", cmdTypeName)
	metadata := ddd_context.GetMetadataContext(ctx)
	return CallMethod(aggregate, methodName, ctx, cmd, metadata)
}
func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...LoadAggregateOption) error {
	rootId := cmd.GetAggregateId().RootId()
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), rootId, aggregate, opts...)
	if err != nil {
		return err
	}
	if !find {
		return ddd_errors.NewAggregateIdNotFondError(rootId)
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
		applog.Error("", "ddd", "NewDomainEvent", err.Error())
		return err
	}
	if err = callEventHandler(ctx, handler, record.EventType, record.EventRevision, event); err != nil {
		applog.Error("", "ddd", "CallEventHandler", err.Error())
	}
	return err
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
