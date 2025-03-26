package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"time"
)

func (d *DaoBase[T]) PublishEvent(ctx context.Context, opeType idao.AccessType, entity T, opts ...*idao.CallOptions) {
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
	case idao.AccessTypeCreate:
		dbevent.CreateEvent(ctx, agg, event)
	case idao.AccessTypeUpdate:
		dbevent.ApplyEvent(ctx, agg, event)
	case idao.AccessTypeDelete:
		dbevent.ApplyEvent(ctx, agg, event)
	}
}

func (d *DaoBase[T]) PublishBatchEvent(ctx context.Context, opeType idao.AccessType, list []map[string]any, opts ...*idao.CallOptions) {
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

func (d *DaoBase[T]) NewEvent(ctx context.Context, operateType idao.AccessType, entity T, opt *idao.CallOptions) (*dbevent.Event, error) {
	o := idao.NewCallOptions(opt)
	eventId := idutils.NewId()
	tenantId := d.store.GetTenantId(entity)
	aggId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	eventType := d.GetEventType(operateType, opt)

	event := dbevent.NewEvent()
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

func (d *DaoBase[T]) NewAggregateAndEvent(ctx context.Context, operateType idao.AccessType, entity T, opts ...*idao.CallOptions) (*dbevent.Aggregate, *dbevent.Event, error) {
	opt := idao.NewCallOptions(opts...)
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

func (d *DaoBase[T]) GetEventType(accessType idao.AccessType, opts *idao.CallOptions) string {
	eventType := d.tableName
	if opts != nil && opts.EventType != nil {
		eventType = *opts.EventType
	}
	return fmt.Sprintf("%s.%s.%s", d.appId, eventType, string(accessType))
}

func (d *DaoBase[T]) NewAggregate(entity T, opt *idao.CallOptions) (*dbevent.Aggregate, error) {
	tenantId := d.store.GetTenantId(entity)
	aggId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	agg := dbevent.NewAggregate()
	agg.TenantId = tenantId
	agg.AggId = aggId
	agg.AggVer = "v1.0"
	agg.AggType = d.tableName
	return agg, nil
}
