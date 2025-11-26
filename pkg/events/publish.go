package events

import (
	"context"
)

func Publish(ctx context.Context, publisher EventPublisher, data IEvent) error {
	event, err := newEventEntity(ctx, data)
	if err != nil {
		return err
	}
	return publisher.Publish(ctx, event)
}

func PublishList[T IEvent](ctx context.Context, publisher EventPublisher, data []T) error {
	var events []*EventEntity
	for _, data := range data {
		event, err := newEventEntity(ctx, data)
		if err != nil {
			return err
		}
		if event == nil {
			continue
		}
		if err = event.Verify(); err != nil {
			return err
		}
		events = append(events, event)
	}
	return publisher.PublishList(ctx, events)
}

func newEventEntity(ctx context.Context, e IEvent) (*EventEntity, error) {
	dataMap, err := anyToMap(e.GetData())
	if err != nil {
		return nil, err
	}
	return &EventEntity{
		Id:          e.GetId(),
		AppId:       e.GetAppId(),
		TenantId:    e.GetTenantId(),
		EventType:   e.GetEventType(),
		Data:        dataMap,
		Meta:        e.GetMeta(),
		CreatedTime: e.GetCreatedTime(),
	}, nil
}

func anyToMap(data any) (map[string]any, error) {
	var dataMap map[string]any
	if val, ok := data.(map[string]any); ok {
		dataMap = val
	} else {
		mapVal, e1 := StructToMap(data)
		if e1 == nil {
			dataMap = mapVal
		} else {
			return nil, e1
		}
	}
	return dataMap, nil
}
