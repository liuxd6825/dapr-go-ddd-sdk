package events

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

func NewOutboxEvent(ctx context.Context, appId string, tenantId string, data any, meta map[string]any) (*dbevent.OutboxEvent, error) {
	var err error
	_, eventType, _ := reflectutils.GetTypeDetails(data)
	eventType = stringutils.MidlineString(eventType)

	meta, err = newMeta(ctx, tenantId, meta)
	if err != nil {
		return nil, err
	}

	dataMap, err := StructToMap(data)
	if err != nil {
		return nil, err
	}
	event := &dbevent.OutboxEvent{
		Id:          idutils.NewUlid2(),
		AppId:       appId,
		TenantId:    tenantId,
		EventType:   eventType,
		Data:        dataMap,
		Meta:        meta,
		CreatedTime: times.NewTime(),
	}
	return event, err
}
