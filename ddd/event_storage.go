package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"strings"
)

var strEmpty = ""

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error)
	GetEvents(ctx context.Context, req *daprclient.GetEventsRequest) (*daprclient.GetEventsResponse, error)
	ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventResponse, error)
	CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error)
	DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error)
	SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error)
	GetRelations(ctx context.Context, req *daprclient.GetRelationsRequest) (*daprclient.GetRelationsResponse, error)
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
	if strings.HasPrefix(ver, "v") || strings.HasPrefix(ver, "V") {
		ver = "V" + ver[1:]
	} else {
		ver = "V" + ver
	}
	return fmt.Sprintf("On%s%s", name, ver)
}
