package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/impl/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"time"
)

func (d *DaoBase) PublishEvent(ctx context.Context, opeType db.AccessType, entity map[string]any, opts ...*db.CallOptions) {
	if !d.isPubEvent {
		return
	}

	agg, event, err := d.NewAggregateAndEvent(ctx, opeType, entity, opts...)
	if err != nil {
		panic(err)
	}
	logs.Debug(ctx, "", logs.Fields{
		"eventId":     event.EventId,
		"eventType":   event.EventType,
		"commandId":   event.CommandId,
		"aggregateId": event.AggregateId,
		"tenantId":    d.GetTenantId(ctx),
	})
	switch opeType {
	case db.AccessTypeCreate:
		events.CreateEvent(ctx, agg, event)
	case db.AccessTypeUpdate:
		events.ApplyEvent(ctx, agg, event)
	case db.AccessTypeDelete:
		events.ApplyEvent(ctx, agg, event)
	}
}

func (d *DaoBase) PublishBatchEvent(ctx context.Context, opeType db.AccessType, list []map[string]any, opts ...*db.CallOptions) {
	if !d.isPubEvent {
		return
	}

	/*
		agg, event, err := d.NewAggregateAndEvent(opeType, list, opts...)
		if err != nil {
			panic(err)
		}
		logs.Debug(ctx, "", logs.Fields{
			"eventId":     event.EventId,
			"eventType":   event.EventType,
			"commandId":   event.CommandId,
			"aggregateId": event.AggregateId,
			"tenantId":    d.GetTenantId(ctx),
		})
		switch opeType {
		case db.AccessTypeCreate:
			server.GetEventPkg().CreateEvent(ctx, agg, event)
		case db.AccessTypeUpdate:
			server.GetEventPkg().ApplyEvent(ctx, agg, event)
		case db.AccessTypeDelete:
			server.GetEventPkg().ApplyEvent(ctx, agg, event)
		}

	*/
}

func (d *DaoBase) NewEvent(ctx context.Context, operateType db.AccessType, entity map[string]any, opt *db.CallOptions) (*common.Event, error) {
	o := db.NewCallOptions(opt)
	eventId := idutils.NewId()
	tenantId := d.dao.GetTenantId(entity)
	aggId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	eventType := d.GetEventType(operateType, opt)

	event := common.NewEvent()
	event.CommandId = o.GetCommandId(idutils.NewId())
	event.EventId = eventId
	event.EventType = eventType
	event.TenantId = tenantId
	event.CreatedTime = time.Now()
	event.AggregateId = aggId
	event.Data = entity
	event.EventVersion = o.GetEventVersion("v1.0")

	return event, nil
}

func (d *DaoBase) NewAggregateAndEvent(ctx context.Context, operateType db.AccessType, entity map[string]any, opts ...*db.CallOptions) (*events.Aggregate, *common.Event, error) {
	opt := db.NewCallOptions(opts...)
	event, err := d.NewEvent(ctx, operateType, entity, opt)
	if err != nil {
		return nil, nil, err
	}
	agg, err := d.NewAggregate(entity, opt)
	if err != nil {
		return nil, nil, err
	}
	return agg, event, nil
}

func (d *DaoBase) GetEventType(accessType db.AccessType, opts *db.CallOptions) string {
	eventType := d.tableName
	if opts != nil && opts.EventType != nil {
		eventType = *opts.EventType
	}
	return common.GetEventType(d.appId, eventType, string(accessType))
}

func (d *DaoBase) NewAggregate(entity map[string]any, opt *db.CallOptions) (*events.Aggregate, error) {
	tenantId := d.dao.GetTenantId(entity)
	aggId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	agg := events.NewAggregate()
	agg.TenantId = tenantId
	agg.AggId = aggId
	agg.AggVer = "v1.0"
	agg.AggType = d.tableName
	return agg, nil
}
