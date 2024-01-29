package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"reflect"
	"strings"
)

var strEmpty = ""

type EventStore interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (Aggregate, bool, error)
	LoadEvent(ctx context.Context, req *dapr.LoadEventsRequest) (*dapr.LoadEventsResponse, error)
	GetEvents(ctx context.Context, req *dapr.GetEventsRequest) (*dapr.GetEventsResponse, error)
	ApplyEvent(ctx context.Context, req *dapr.ApplyEventRequest) (*dapr.ApplyEventResponse, error)
	Commit(ctx context.Context, req *dapr.CommitRequest) (res *dapr.CommitResponse, resErr error)
	Rollback(ctx context.Context, req *dapr.RollbackRequest) (res *dapr.RollbackResponse, resErr error)
	SaveSnapshot(ctx context.Context, req *dapr.SaveSnapshotRequest) (*dapr.SaveSnapshotResponse, error)
	GetRelations(ctx context.Context, req *dapr.GetRelationsRequest) (*dapr.GetRelationsResponse, error)
	GetPubsubName() string
}

var snapshotEventsMinCount = 20

type CallEventType int

const (
	EventCreate CallEventType = iota
	EventApply
	EventDelete
)

type EventStoreOption func(EventStore)

func (t CallEventType) ToString() string {
	switch t {
	case EventCreate:
		return "create"
	case EventApply:
		return "apply"
	case EventDelete:
		return "delete"
	}
	return ""
}

func PubsubName(pubsubName string) EventStoreOption {
	return func(es EventStore) {
		s, _ := es.(*grpcEventStore)
		s.pubsubName = pubsubName
	}
}

func checkEvent(aggregate Aggregate, event DomainEvent) error {
	if err := assert.NotNil(event, assert.NewOptions("event is nil")); err != nil {
		return err
	}
	if err := assert.NotNil(aggregate, assert.NewOptions("agg is nil")); err != nil {
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

// CallEventHandler
// @Description: 调用领域事件监听器
// @param ctx
// @param handler
// @param record
// @return error
func CallEventHandler(ctx context.Context, handler interface{}, record *dapr.EventRecord) error {
	if record == nil {
		return errors.New("package:ddd; func:CallEventHandler(); error: record is nil")
	}
	event, err := NewDomainEvent(record)
	if err != nil {
		return errors.New("package:ddd; func:CallEventHandler(); error: NewDomainEvent() %v", err.Error())
	}
	metadata := record.Metadata
	return callEventHandler(ctx, handler, record.EventType, record.EventVersion, event, metadata)
}

// callEventHandler
//
//	@Description: 调用QueryHandler事件订阅处理器
//	@param ctx  上下文
//	@param queryHandler 事件处理器
//	@param eventType  事件类型
//	@param eventVersion 事件版本号
//	@param event 事件对象
//	@param metadata  事件元数据
//	@return error 错误
func callEventHandler(ctx context.Context, queryHandler any, eventType string, eventVersion string, event any, metadata Metadata) error {
	methodName := getEventMethodName(eventType, eventVersion)
	return callMethod(ctx, queryHandler, methodName, event, metadata)
}

// callCommandHandler
//
//	@Description: 调用CommandHandle命令处理器
//	@param ctx 上下文
//	@param aggregate 聚合根对象
//	@param cmd 命令
//	@return error 错误
func callCommandHandler(ctx context.Context, aggregate any, cmd Command) error {
	cmdTypeName := reflect.ValueOf(cmd).Elem().Type().Name()
	methodName := fmt.Sprintf("%s", cmdTypeName)
	metadata := ddd_context.GetMetadataContext(ctx)
	return callMethod(ctx, aggregate, methodName, cmd, metadata)
}

// callMethod
//
//	@Description: 动态调用方法
//	@param obj 方法对象
//	@param methodName 方法名称
//	@param ctx
//	@param eventOrCommand
//	@param metadata
//	@return error
func callMethod(ctx context.Context, obj any, methodName string, eventOrCommand any, metadata Metadata) error {
	return reflectutils.CallMethod(obj, methodName, ctx, eventOrCommand, metadata)
}

// getEventMethodName
// @Description: 根据事件类型名称获取接受事件方法名称
// @param eventType
// @param revision
// @return string
func getEventMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "s", -1)
	if strings.HasPrefix(ver, "v") || strings.HasPrefix(ver, "V") {
		ver = "V" + ver[1:]
	} else {
		ver = "V" + ver
	}
	return fmt.Sprintf("On%s%s", name, ver)
}
