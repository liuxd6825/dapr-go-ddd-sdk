package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type ApplyEventOptions struct {
	pubsubName       *string
	metadata         *map[string]string
	eventStorageKey  *string
	sessionId        *string
	closeEventSource *bool
}

func NewApplyEventOptions(metadata *map[string]string) *ApplyEventOptions {
	return &ApplyEventOptions{
		metadata: metadata,
	}
}

func NewApplyEventOptionsNil() *ApplyEventOptions {
	return &ApplyEventOptions{}
}

func OptionCloseEventSource() *ApplyEventOptions {
	t := true
	return &ApplyEventOptions{closeEventSource: &t}
}

func (a *ApplyEventOptions) Merge(opts ...*ApplyEventOptions) *ApplyEventOptions {
	for _, opt := range opts {
		if opt.eventStorageKey != nil {
			a.eventStorageKey = opt.eventStorageKey
		}
		if opt.metadata != nil {
			a.metadata = opt.metadata
		}
		if opt.pubsubName != nil {
			a.pubsubName = opt.pubsubName
		}
		if opt.sessionId != nil {
			a.sessionId = opt.sessionId
		}
		if opt.closeEventSource != nil {
			a.closeEventSource = opt.closeEventSource
		}
	}
	return a
}

func (a *ApplyEventOptions) SetCloseEventSource(v bool) *ApplyEventOptions {
	a.closeEventSource = &v
	return a
}

func (a *ApplyEventOptions) GetCloseEventSource() *bool {
	return a.closeEventSource
}

func (a *ApplyEventOptions) SetPubsubName(pubsubName string) *ApplyEventOptions {
	a.pubsubName = &pubsubName
	return a
}

func (a *ApplyEventOptions) GetPubsubName() *string {
	return a.pubsubName
}

func (a *ApplyEventOptions) SetEventStorageKey(eventStorageKey string) *ApplyEventOptions {
	a.eventStorageKey = &eventStorageKey
	return a
}

func (a *ApplyEventOptions) GetEventStorageKey() string {
	if a.eventStorageKey != nil {
		return *a.eventStorageKey
	}
	return ""
}

func (a *ApplyEventOptions) SetMetadata(value *map[string]string) *ApplyEventOptions {
	a.metadata = value
	return a
}

func (a *ApplyEventOptions) GetMetadata() *map[string]string {
	return a.metadata
}

func (a *ApplyEventOptions) SetSessionId(value string) *ApplyEventOptions {
	a.sessionId = &value
	return a
}

func (a *ApplyEventOptions) GetSessionId() *string {
	return a.sessionId
}

func ApplyEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.ApplyEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventApply, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*daprclient.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func ApplyEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*daprclient.ApplyEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventApply, aggregate, events, opts...)
	if resp, ok := res.(*daprclient.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.CreateEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventCreate, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*daprclient.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*daprclient.CreateEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventCreate, aggregate, events, opts...)
	if resp, ok := res.(*daprclient.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*daprclient.DeleteEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventDelete, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*daprclient.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*daprclient.DeleteEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventDelete, aggregate, events, opts...)
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
func callDaprEventMethod(ctx context.Context, callEventType CallEventType, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (resAny any, resErr error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			resErr = e
		}
	}()

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(&metadata).Merge(opts...)

	for _, event := range events {
		if err := checkEvent(aggregate, event); err != nil {
			return nil, err
		}
	}

	tenantId := aggregate.GetTenantId()
	aggId := aggregate.GetAggregateId()
	aggType := aggregate.GetAggregateType()

	sessionId := ""
	if options.GetSessionId() != nil {
		sessionId = *options.GetSessionId()
	}

	session, ok := getSession(ctx)
	if ok && session != nil {
		sessionId = session.sessionId
	}

	logInfo := &applog.LogInfo{
		TenantId:  aggregate.GetTenantId(),
		ClassName: "ddd",
		FuncName:  "callDaprEventMethod",
		Message:   fmt.Sprintf("%v", aggregate),
		Level:     logs.InfoLevel,
	}

	var err error
	var res any
	err = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		var eventStorage EventStorage
		applyEvents := make([]*daprclient.EventDto, 0)

		eventStorage, err = GetEventStorage(options.GetEventStorageKey())
		if err != nil {
			return nil, err
		}

		defaultIsSourcing := true
		// 判断是否需要进行"事件溯源"控制
		if options.closeEventSource != nil {
			closeEs := *options.closeEventSource
			defaultIsSourcing = !closeEs
		}

		for _, event := range events {
			relation, _, err := GetRelationByStructure(event.GetData())
			if err != nil {
				return nil, err
			}
			isSourcing := defaultIsSourcing
			if e, ok := event.(IsSourcing); ok {
				isSourcing = e.GetIsSourcing()
			}
			eventDto := &daprclient.EventDto{
				ApplyType:    callEventType.ToString(),
				CommandId:    event.GetCommandId(),
				EventId:      event.GetEventId(),
				EventVersion: event.GetEventVersion(),
				EventType:    event.GetEventType(),
				Metadata:     *options.metadata,
				PubsubName:   eventStorage.GetPubsubName(),
				EventData:    event,
				Relations:    relation,
				Topic:        event.GetEventType(),
				IsSourcing:   isSourcing,
			}
			applyEvents = append(applyEvents, eventDto)
		}

		res, err = applyEvent(ctx, eventStorage, tenantId, sessionId, aggId, aggType, applyEvents)
		if err != nil {
			return nil, err
		}

		if defaultIsSourcing {
			for _, event := range events {
				if err = callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventVersion(), event); err != nil {
					return nil, err
				}
			}
		}
		return res, nil
	})

	go func() {
		_ = callActorSaveSnapshot(ctx, tenantId, aggId, aggType)
	}()

	return res, err
}

func applyEvent(ctx context.Context, eventStorage EventStorage, tenantId, sessionId, aggregateId, aggregateType string, events []*daprclient.EventDto) (*daprclient.ApplyEventResponse, error) {
	req := &daprclient.ApplyEventRequest{
		SessionId:     sessionId,
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.ApplyEvent(ctx, req)
	return resp, err
}

func Commit(ctx context.Context, tenantId string, sessionId string, opts ...*ApplyEventOptions) (res *daprclient.CommitResponse, resErr error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			resErr = e
		}
	}()
	req := &daprclient.CommitRequest{
		TenantId:  tenantId,
		SessionId: sessionId,
	}

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(&metadata).Merge(opts...)

	eventStorage, err := GetEventStorage(options.GetEventStorageKey())
	if err != nil {
		return nil, err
	}

	out, err := eventStorage.Commit(ctx, req)
	resp := &daprclient.CommitResponse{
		Headers: out.Headers,
	}
	return resp, err
}

func Rollback(ctx context.Context, tenantId string, sessionId string, opts ...*ApplyEventOptions) (res *daprclient.RollbackResponse, resErr error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			resErr = e
		}
	}()
	req := &daprclient.RollbackRequest{
		TenantId:  tenantId,
		SessionId: sessionId,
	}

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(&metadata).Merge(opts...)

	eventStorage, err := GetEventStorage(options.GetEventStorageKey())
	if err != nil {
		return nil, err
	}

	out, err := eventStorage.Rollback(ctx, req)
	resp := &daprclient.RollbackResponse{
		Headers: out.Headers,
	}
	return resp, err
}

/*func createEvent(ctx context.Context, sessionId string, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) (*daprclient.CreateEventResponse, error) {
	req := &daprclient.CreateEventRequest{

		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.CreateEvent(ctx, req)
	return resp, err
}

func deleteEvent(ctx context.Context, sessionId string, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, event *daprclient.EventDto) (*daprclient.DeleteEventResponse, error) {
	req := &daprclient.DeleteEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Event:         event,
	}
	resp, err := eventStorage.DeleteEvent(ctx, req)
	return resp, err
}*/

//
//  callActorSaveSnapshot
//  @Description: 通过调用 actor service 生成聚合快照。
//  @param ctx
//  @param tenantId
//  @param aggregateId
//  @param aggregateType
//  @return error
//
func callActorSaveSnapshot(ctx context.Context, tenantId, aggregateId, aggregateType string) (resErr error) {
	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			resErr = e
		}
	}()

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
