package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type Metadata = map[string][]string

type ApplyEventOptions struct {
	PubsubName       *string
	EventStoreName   *string
	Metadata         Metadata
	SessionId        *string
	CloseEventSource *bool
}

func NewApplyEventOptions(metadata Metadata) *ApplyEventOptions {
	return &ApplyEventOptions{
		Metadata: metadata,
	}
}

func NewApplyEventOptionsNil() *ApplyEventOptions {
	return &ApplyEventOptions{}
}

func OptionCloseEventSource() *ApplyEventOptions {
	t := true
	return &ApplyEventOptions{CloseEventSource: &t}
}

func (a *ApplyEventOptions) Merge(opts ...*ApplyEventOptions) *ApplyEventOptions {
	for _, opt := range opts {
		if opt.EventStoreName != nil {
			a.EventStoreName = opt.EventStoreName
		}
		if opt.Metadata != nil {
			a.Metadata = opt.Metadata
		}
		if opt.PubsubName != nil {
			a.PubsubName = opt.PubsubName
		}
		if opt.SessionId != nil {
			a.SessionId = opt.SessionId
		}
		if opt.CloseEventSource != nil {
			a.CloseEventSource = opt.CloseEventSource
		}
	}
	return a
}

func (a *ApplyEventOptions) SetMetadataFromCtx(ctx context.Context) *ApplyEventOptions {
	if ctx == nil {
		return a
	}
	if header, ok := appctx.GetHeader(ctx); ok {
		a.Metadata = header
	}
	if tenantId, ok := appctx.GetTenantId(ctx); ok {
		a.Metadata["TenantId"] = []string{tenantId}
	}
	if auth, ok := appctx.GetAuthToken(ctx); ok {
		a.Metadata[Authorization] = []string{auth.GetToken()}
	}
	return a
}

func (a *ApplyEventOptions) SetCloseEventSource(v bool) *ApplyEventOptions {
	a.CloseEventSource = &v
	return a
}

func (a *ApplyEventOptions) GetCloseEventSource() *bool {
	return a.CloseEventSource
}

func (a *ApplyEventOptions) SetPubsubName(pubsubName string) *ApplyEventOptions {
	a.PubsubName = &pubsubName
	return a
}

func (a *ApplyEventOptions) GetPubsubName() *string {
	return a.PubsubName
}

func (a *ApplyEventOptions) SetEventStoreKey(eventStoreName string) *ApplyEventOptions {
	a.EventStoreName = &eventStoreName
	return a
}

func (a *ApplyEventOptions) GetEventStoreName() string {
	if a.EventStoreName != nil {
		return *a.EventStoreName
	}
	return ""
}

func (a *ApplyEventOptions) SetMetadata(value Metadata) *ApplyEventOptions {
	a.Metadata = value
	return a
}

func (a *ApplyEventOptions) GetMetadata() Metadata {
	if a.Metadata == nil {
		a.Metadata = Metadata{}
	}
	return a.Metadata
}

func (a *ApplyEventOptions) SetSessionId(value string) *ApplyEventOptions {
	a.SessionId = &value
	return a
}

func (a *ApplyEventOptions) GetSessionId() *string {
	return a.SessionId
}

func ApplyEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	res, err := publishEvents(ctx, EventApply, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func ApplyEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	res, err := publishEvents(ctx, EventApply, aggregate, events, opts...)
	if resp, ok := res.(*dapr.ApplyEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	res, err := publishEvents(ctx, EventCreate, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func CreateEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	res, err := publishEvents(ctx, EventCreate, aggregate, events, opts...)
	if resp, ok := res.(*dapr.CreateEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	res, err := publishEvents(ctx, EventDelete, aggregate, []DomainEvent{event}, opts...)
	if resp, ok := res.(*dapr.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

func DeleteEvents(ctx context.Context, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	res, err := publishEvents(ctx, EventDelete, aggregate, events, opts...)
	if resp, ok := res.(*dapr.DeleteEventResponse); ok {
		return resp, err
	}
	return nil, err
}

// callDaprEventMethod
// @Description: 应用领域事件
// @param ctx
// @param agg
// @param event
// @param options
// @return err
func publishEvents(ctx context.Context, callEventType CallEventType, aggregate Aggregate, events []DomainEvent, opts ...*ApplyEventOptions) (resAny any, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	options := NewApplyEventOptionsNil().SetMetadataFromCtx(ctx).Merge(opts...)

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
		errs.AddString("callDaprEventMethod() aggregate.GetTenantId() is empty")
	}
	if len(aggId) == 0 {
		errs.AddString("callDaprEventMethod()  aggregate.GetAggregateId() is empty")
	}
	if len(aggType) == 0 {
		errs.AddString("callDaprEventMethod() aggregate.GetAggregateType() is empty")
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
		if options.CloseEventSource != nil {
			closeEs := *options.CloseEventSource
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
				Metadata:     options.GetMetadata(),
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
				if err = callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventVersion(), event, options.GetMetadata()); err != nil {
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

	options := NewApplyEventOptionsNil().SetMetadataFromCtx(ctx).Merge(opts...)

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

	options := NewApplyEventOptionsNil().SetMetadataFromCtx(ctx).Merge(opts...)

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

/*func createEvent(ctx context.Context, SessionId string, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) (*daprclient.CreateEventResponse, error) {
	req := &daprclient.CreateEventRequest{

		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	resp, err := eventStorage.CreateEvent(ctx, req)
	return resp, err
}

func deleteEvent(ctx context.Context, SessionId string, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, event *daprclient.EventDto) (*daprclient.DeleteEventResponse, error) {
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
