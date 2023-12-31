package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type ApplyEventOptions struct {
	pubsubName      *string
	metadata        *map[string]string
	eventStorageKey *string
}

func ApplyEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.ApplyEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventApply, aggregate, event, opts...)
	if resp, ok := res.(*daprclient.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.CreateEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventCreate, aggregate, event, opts...)
	if resp, ok := res.(*daprclient.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.DeleteEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventDelete, aggregate, event, opts...)
	if resp, ok := res.(*daprclient.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

//
// callDaprEventMethod
// @Description: 应用领域事件
// @param ctx
// @param aggregate
// @param event
// @param options
// @return err
//
func callDaprEventMethod(ctx context.Context, callEventType CallEventType, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (any, error) {
	if err := checkEvent(aggregate, event); err != nil {
		return nil, err
	}

	tenantId := event.GetTenantId()
	aggId := event.GetAggregateId()
	aggType := aggregate.GetAggregateType()

	metadata := make(map[string]string)
	options := &ApplyEventOptions{
		pubsubName:      &strEmpty,
		metadata:        &metadata,
		eventStorageKey: &strEmpty,
	}
	for _, opt := range opts {
		if opt.eventStorageKey != nil {
			options.eventStorageKey = opt.eventStorageKey
		}
		if opt.metadata != nil {
			options.metadata = opt.metadata
		}
		if opt.pubsubName != nil {
			options.pubsubName = opt.pubsubName
		}
	}

	logInfo := &applog.LogInfo{
		TenantId:  aggregate.GetTenantId(),
		ClassName: "ddd",
		FuncName:  "callDaprEventMethod",
		Message:   fmt.Sprintf("%v", aggregate),
		Level:     applog.INFO,
	}

	var err error
	var res any
	err = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		var eventStorage EventStorage
		eventStorage, err = GetEventStorage(*options.eventStorageKey)
		if err != nil {
			return nil, err
		}
		relation, _, err := GetRelationByStructure(event.GetData())
		if err != nil {
			return nil, err
		}
		applyEvents := []*daprclient.EventDto{
			{
				CommandId:    event.GetCommandId(),
				EventId:      event.GetEventId(),
				EventVersion: event.GetEventVersion(),
				EventType:    event.GetEventType(),
				Metadata:     *options.metadata,
				PubsubName:   *options.pubsubName,
				EventData:    event,
				Relations:    relation,
				Topic:        event.GetEventType(),
			},
		}
		if callEventType == EventCreate {
			res, err = createEvent(ctx, eventStorage, tenantId, aggId, aggType, applyEvents)
		} else if callEventType == EventApply {
			res, err = applyEvent(ctx, eventStorage, tenantId, aggId, aggType, applyEvents)
		} else if callEventType == EventDelete {
			res, err = deleteEvent(ctx, eventStorage, tenantId, aggId, aggType, applyEvents[0])
		}
		if err != nil {
			return nil, err
		}
		if err = callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventVersion(), event); err != nil {
			return nil, err
		}
		return res, nil
	})

	go func() {
		_ = callActorSaveSnapshot(ctx, tenantId, aggId, aggType)
	}()

	return res, err
}

func applyEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) (*daprclient.ApplyEventResponse, error) {
	req := &daprclient.ApplyEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.ApplyEvent(ctx, req)
	return resp, err
}

func createEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) (*daprclient.CreateEventResponse, error) {
	req := &daprclient.CreateEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.CreateEvent(ctx, req)
	return resp, err
}

func deleteEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, event *daprclient.EventDto) (*daprclient.DeleteEventResponse, error) {
	req := &daprclient.DeleteEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Event:         event,
	}
	resp, err := eventStorage.DeleteEvent(ctx, req)
	return resp, err
}

//
//  callActorSaveSnapshot
//  @Description: 通过调用 actor service 生成聚合快照。
//  @param ctx
//  @param tenantId
//  @param aggregateId
//  @param aggregateType
//  @return error
//
func callActorSaveSnapshot(ctx context.Context, tenantId, aggregateId, aggregateType string) error {
	client, err := daprclient.GetDaprDDDClient().DaprClient()
	if err != nil {
		return err
	}
	snapshotClient := NewAggregateSnapshotClient(client, aggregateType, aggregateId)
	_, err = snapshotClient.SaveSnapshot(ctx, &SaveSnapshotRequest{
		TenantId:      tenantId,
		AggregateType: aggregateType,
		AggregateId:   aggregateId,
	})
	return err
}

func (a *ApplyEventOptions) SetPubsubName(pubsubName string) *ApplyEventOptions {
	a.pubsubName = &pubsubName
	return a
}

func (a *ApplyEventOptions) SetEventStorageKey(eventStorageKey string) *ApplyEventOptions {
	a.eventStorageKey = &eventStorageKey
	return a
}

func (a *ApplyEventOptions) SetMetadata(value *map[string]string) *ApplyEventOptions {
	a.metadata = value
	return a
}

func NewApplyEventOptions(metadata *map[string]string) *ApplyEventOptions {
	return &ApplyEventOptions{
		metadata: metadata,
	}
}
