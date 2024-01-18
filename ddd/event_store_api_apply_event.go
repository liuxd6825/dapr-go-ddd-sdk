package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type ApplyEventOptions struct {
	pubsubName       *string
	eventStoreName   *string
	metadata         map[string]string
	sessionId        *string
	closeEventSource *bool
}

func NewApplyEventOptions(metadata map[string]string) *ApplyEventOptions {
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
		if opt.eventStoreName != nil {
			a.eventStoreName = opt.eventStoreName
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

func (a *ApplyEventOptions) SetEventStoreKey(eventStoreName string) *ApplyEventOptions {
	a.eventStoreName = &eventStoreName
	return a
}

func (a *ApplyEventOptions) GetEventStoreName() string {
	if a.eventStoreName != nil {
		return *a.eventStoreName
	}
	return ""
}

func (a *ApplyEventOptions) SetMetadata(value map[string]string) *ApplyEventOptions {
	a.metadata = value
	return a
}

func (a *ApplyEventOptions) GetMetadata() map[string]string {
	return a.metadata
}

func (a *ApplyEventOptions) SetSessionId(value string) *ApplyEventOptions {
	a.sessionId = &value
	return a
}

func (a *ApplyEventOptions) GetSessionId() *string {
	return a.sessionId
}

func ApplyEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventApply, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func ApplyEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventApply, aggregate, events, opts...)
	if resp, ok := res.(*dapr.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventCreate, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventCreate, aggregate, events, opts...)
	if resp, ok := res.(*dapr.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventDelete, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	res, err := callDaprEventMethod(ctx, EventDelete, aggregate, events, opts...)
	if resp, ok := res.(*dapr.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

// callDaprEventMethod
// @Description: 应用领域事件
// @param ctx
// @param aggregate
// @param event
// @param options
// @return err
func callDaprEventMethod(ctx context.Context, callEventType CallEventType, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (resAny any, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(metadata).Merge(opts...)

	for _, event := range events {
		if err := checkEvent(aggregate, event); err != nil {
			return nil, err
		}
	}

	tenantId := aggregate.GetTenantId()
	aggId := aggregate.GetAggregateId()
	aggType := aggregate.GetAggregateType()

	errs := errors.NewErrors()
	if len(tenantId) == 0 {
		errs.AddString("tenantId is empty")
	}
	if len(aggId) == 0 {
		errs.AddString("aggregateId is empty")
	}
	if len(aggType) == 0 {
		errs.AddString("aggregateType is empty")
	}
	if !errs.IsEmpty() {
		return nil, errs
	}

	sessionId := ""
	if options.GetSessionId() != nil {
		sessionId = *options.GetSessionId()
	}

	session, ok := getSession(ctx)
	if ok && session != nil {
		sessionId = session.sessionId
	}

	field := logs.Fields{
		"package":       "ddd",
		"funcName":      "callDaprEventMethod",
		"aggregateId":   aggId,
		"aggregateType": aggType,
	}

	var err error
	var res any

	//默认事件溯源为true
	defaultIsSourcing := true

	err = logs.DebugStart(ctx, tenantId, field, func() error {
		var eventStore EventStore
		applyEvents := make([]*dapr.EventDto, 0)

		eventStore, err = GetEventStore(options.GetEventStoreName())
		if err != nil {
			return err
		}

		pubsubName := eventStore.GetPubsubName()
		if val := options.GetPubsubName(); val != nil {
			pubsubName = *val
		}

		// 判断是否需要进行"事件溯源"控制
		if options.closeEventSource != nil {
			closeEs := *options.closeEventSource
			defaultIsSourcing = !closeEs
		}

		for _, event := range events {
			relation, _, err := GetRelationByStructure(event.GetData())
			if err != nil {
				return err
			}
			isSourcing := defaultIsSourcing
			if e, ok := event.(IsSourcing); ok {
				isSourcing = e.GetIsSourcing()
			}
			eventDto := &dapr.EventDto{
				ApplyType:    callEventType.ToString(),
				CommandId:    event.GetCommandId(),
				EventId:      event.GetEventId(),
				EventVersion: event.GetEventVersion(),
				EventType:    event.GetEventType(),
				Metadata:     options.metadata,
				PubsubName:   pubsubName,
				EventData:    event,
				Relations:    relation,
				Topic:        event.GetEventType(),
				IsSourcing:   isSourcing,
			}
			applyEvents = append(applyEvents, eventDto)
		}

		res, err = applyEvent(ctx, eventStore, tenantId, sessionId, aggId, aggType, applyEvents)
		if err != nil {
			return err
		}

		if defaultIsSourcing {
			for _, event := range events {
				if err = callEventHandler(ctx, aggregate, event.GetTenantId(), event.GetEventType(), event.GetEventVersion(), event.GetEventId(), event); err != nil {
					return err
				}
			}
		}

		// 如果是溯源模式进行聚合根镜像
		if defaultIsSourcing && len(events) > 100 {
			go func() {
				_ = callActorSaveSnapshot(ctx, tenantId, aggId, aggType)
			}()
		}

		return nil
	})

	return res, err
}

func applyEvent(ctx context.Context, eventStorage EventStore, tenantId, sessionId, aggregateId, aggregateType string, events []*dapr.EventDto) (*dapr.ApplyEventResponse, error) {
	req := &dapr.ApplyEventRequest{
		SessionId:     sessionId,
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.ApplyEvent(ctx, req)
	return resp, err
}

func Commit(ctx context.Context, tenantId string, sessionId string, opts ...*ApplyEventOptions) (res *dapr.CommitResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req := &dapr.CommitRequest{
		TenantId:  tenantId,
		SessionId: sessionId,
	}

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(metadata).Merge(opts...)

	eventStorage, err := GetEventStore(options.GetEventStoreName())
	if err != nil {
		return nil, err
	}

	out, err := eventStorage.Commit(ctx, req)
	resp := &dapr.CommitResponse{
		Headers: out.Headers,
	}
	return resp, err
}

func Rollback(ctx context.Context, tenantId string, sessionId string, opts ...*ApplyEventOptions) (res *dapr.RollbackResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req := &dapr.RollbackRequest{
		TenantId:  tenantId,
		SessionId: sessionId,
	}

	metadata := make(map[string]string)
	options := NewApplyEventOptionsNil().SetMetadata(metadata).Merge(opts...)

	eventStorage, err := GetEventStore(options.GetEventStoreName())
	if err != nil {
		return nil, err
	}

	out, err := eventStorage.Rollback(ctx, req)
	resp := &dapr.RollbackResponse{
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

// callActorSaveSnapshot
// @Description: 通过调用 actor service 生成聚合快照。
// @param ctx
// @param tenantId
// @param aggregateId
// @param aggregateType
// @return error
func callActorSaveSnapshot(ctx context.Context, tenantId, aggregateId, aggregateType string) (resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	client, err := dapr.GetDaprClient().Client()
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
