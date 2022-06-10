package ddd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"reflect"
	"strings"
)

var strEmpty = ""

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventResponse, error)
	CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error)
	DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error)
	SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error)
	GetPubsubName() string
}

var snapshotEventsMinCount = 20

type CallEventType int

const (
	EventCreate CallEventType = iota
	EventApply
	EventDelete
)

type EventStorageOption func(EventStorage)

func PubsubName(pubsubName string) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*grpcEventStorage)
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

//
// LoadAggregate
// @Description: 加载聚合根
// @param ctx 上下文
// @param tenantId 租户id
// @param aggregateId 聚合根id
// @param aggregate 聚合根对象
// @param opts 可选参数
// @return agg    聚合根对象
// @return isFound 是否找到
// @return err 错误
//
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

func SaveSnapshot(ctx context.Context, tenantId string, aggregateType string, aggregateId string, eventStorageKey string) error {
	aggregate, err := NewAggregate(aggregateType)
	if err != nil {
		return err
	}

	req := &daprclient.LoadEventsRequest{
		TenantId:    tenantId,
		AggregateId: aggregateId,
	}
	resp, err := LoadEvents(ctx, req, "")
	if err != nil {
		return err
	}
	if resp.Snapshot == nil && (resp.EventRecords == nil || len(*resp.EventRecords) == 0) {
		return err
	}

	if resp.Snapshot != nil {
		bytes, err := json.Marshal(resp.Snapshot.AggregateData)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, aggregate)
		if err != nil {
			return err
		}
	}
	records := *resp.EventRecords
	if records != nil && len(records) > snapshotEventsMinCount {
		sequenceNumber := uint64(0)
		for _, record := range *resp.EventRecords {
			sequenceNumber = record.SequenceNumber
			if err = CallEventHandler(ctx, aggregate, &record); err != nil {
				return err
			}
		}

		snapshot := &daprclient.SaveSnapshotRequest{
			TenantId:         tenantId,
			AggregateData:    aggregate,
			AggregateId:      aggregateId,
			AggregateType:    aggregateType,
			AggregateVersion: aggregate.GetAggregateVersion(),
			SequenceNumber:   sequenceNumber,
		}
		eventStorage, err := GetEventStorage(eventStorageKey)
		if err != nil {
			return err
		}
		_, err = eventStorage.SaveSnapshot(ctx, snapshot)
		if err != nil {
			return err
		}
	}
	return err
}

//
// LoadEvents
// @Description: 获取领域事件
// @param ctx 上下文
// @param req 传入参数
// @param eventStorageKey 事件存储器key
// @return resp 响应体
// @return err 错误
//
func LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest, eventStorageKey string) (resp *daprclient.LoadEventsResponse, err error) {
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
			return nil, err
		}
		resp, err = eventStorage.LoadEvent(ctx, req)
		return resp, err
	})
	return
}

type ApplyEventOptions struct {
	pubsubName      *string
	metadata        *map[string]string
	eventStorageKey *string
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

//
// callDaprEventMethod
// @Description: 应用领域事件
// @param ctx
// @param aggregate
// @param event
// @param options
// @return err
//
func callDaprEventMethod(ctx context.Context, callEventType CallEventType, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (err error) {
	if err := checkEvent(aggregate, event); err != nil {
		return err
	}
	tenantId := event.GetTenantId()
	aggregateId := event.GetAggregateId()
	aggregateType := aggregate.GetAggregateType()

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
	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		eventStorage, e := GetEventStorage(*options.eventStorageKey)
		if e != nil {
			err = e
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
				Topic:        event.GetEventType(),
			},
		}
		err = nil
		if callEventType == EventCreate {
			err = createEvent(ctx, eventStorage, tenantId, aggregateId, aggregateType, applyEvents)
		} else if callEventType == EventApply {
			err = applyEvent(ctx, eventStorage, tenantId, aggregateId, aggregateType, applyEvents)
		} else if callEventType == EventDelete {
			err = deleteEvent(ctx, eventStorage, tenantId, aggregateId, aggregateType, applyEvents[0])
		}
		if err != nil {
			return nil, err
		}
		if err = callEventHandler(ctx, aggregate, event.GetEventType(), event.GetEventVersion(), event); err != nil {
			return nil, err
		}
		return nil, nil
	})

	go func() {
		_ = callActorSaveSnapshot(ctx, tenantId, aggregateId, aggregateType)
	}()

	return
}

func applyEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) error {
	req := &daprclient.ApplyEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	_, err := eventStorage.ApplyEvent(ctx, req)
	return err
}

func createEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, events []*daprclient.EventDto) error {
	req := &daprclient.CreateEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Events:        events,
	}
	_, err := eventStorage.CreateEvent(ctx, req)
	return err
}

func deleteEvent(ctx context.Context, eventStorage EventStorage, tenantId, aggregateId, aggregateType string, event *daprclient.EventDto) error {
	req := &daprclient.DeleteEventRequest{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Event:         event,
	}
	_, err := eventStorage.DeleteEvent(ctx, req)
	return err
}

func ApplyEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (err error) {
	return callDaprEventMethod(ctx, EventApply, aggregate, event, opts...)
}

func CreateEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (err error) {
	return callDaprEventMethod(ctx, EventCreate, aggregate, event, opts...)
}

func DeleteEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, opts ...*ApplyEventOptions) (err error) {
	return callDaprEventMethod(ctx, EventDelete, aggregate, event, opts...)
}

func checkEvent(aggregate Aggregate, event DomainEvent) error {
	if err := assert.NotNil(event, assert.NewOptions("event is nil")); err != nil {
		return err
	}
	if err := assert.NotNil(aggregate, assert.NewOptions("aggregate is nil")); err != nil {
		return err
	}

	tenantId := event.GetTenantId()
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return err
	}

	aggId := event.GetAggregateId()
	if err := assert.NotEmpty(aggId, assert.NewOptions("aggregateId is empty")); err != nil {
		return err
	}

	aggregateType := aggregate.GetAggregateType()
	if err := assert.NotEmpty(aggregateType, assert.NewOptions("aggregateType is empty")); err != nil {
		return err
	}
	return nil
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

type CreateAggregateOptions struct {
	eventStorageKey *string
}

func (o *CreateAggregateOptions) SetEventStorageKey(eventStorageKey string) {
	o.eventStorageKey = &eventStorageKey
}

//
// CreateAggregate
// @Description: 创建聚合根
// @param ctx
// @param aggregate
// @param cmd
// @param opts
// @return error
//
func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...*CreateAggregateOptions) error {
	options := &CreateAggregateOptions{
		eventStorageKey: &strEmpty,
	}
	for _, item := range opts {
		if item.eventStorageKey != nil {
			options.eventStorageKey = item.eventStorageKey
		}
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

func callCommandHandler(ctx context.Context, aggregate Aggregate, cmd Command) error {
	cmdTypeName := reflect.ValueOf(cmd).Elem().Type().Name()
	methodName := fmt.Sprintf("%s", cmdTypeName)
	metadata := ddd_context.GetMetadataContext(ctx)
	return CallMethod(aggregate, methodName, ctx, cmd, metadata)
}

//
// CommandAggregate
// @Description: 执行聚合命令
// @param ctx
// @param aggregate
// @param cmd
// @param opts
// @return error
//
func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...LoadAggregateOption) error {
	aggId := cmd.GetAggregateId().RootId()
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), aggId, aggregate, opts...)
	if err != nil {
		return err
	}
	if !find {
		return ddd_errors.NewAggregateIdNotFondError(aggId)
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

//
// CallEventHandler
// @Description: 调用领域事件监听器
// @param ctx
// @param handler
// @param record
// @return error
//
func CallEventHandler(ctx context.Context, handler interface{}, record *daprclient.EventRecord) error {
	event, err := NewDomainEvent(record)
	if err != nil {
		_, _ = applog.Error("", "ddd", "NewDomainEvent", err.Error())
		return err
	}
	if err = callEventHandler(ctx, handler, record.EventType, record.EventVersion, event); err != nil {
		_, _ = applog.Error("", "ddd", "CallEventHandler", err.Error())
	}
	return err
}

func callEventHandler(ctx context.Context, handler interface{}, eventType string, eventRevision string, event interface{}) error {
	methodName := getEventMethodName(eventType, eventRevision)
	return CallMethod(handler, methodName, ctx, event)
}

//
//  getEventMethodName
//  @Description: 根据事件类型名称获取接受事件方法名称
//  @param eventType
//  @param revision
//  @return string
//
func getEventMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "s", -1)
	return fmt.Sprintf("On%sV%s", name, ver)
}
