package events

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"sync"
)

var _outboxDao idao.OutboxEventDao
var _outboxOnce sync.Once

func Publish(ctx context.Context, appId string, data any, meta any) error {
	tenantId, err := appctx.GetTenantId3(ctx)
	if err != nil {
		return err
	}
	outboxDao := newOutboxEventDao()
	event, err := getEvent(ctx, appId, tenantId, data, meta)
	if err != nil {
		return err
	}
	return outboxDao.Create(ctx, event).GetError()
}

func newOutboxEventDao() idao.OutboxEventDao {
	_outboxOnce.Do(func() {
		_outboxDao = dao.NewOutboxEventDao(config.DBKey)
	})
	return _outboxDao
}

func getEvent(ctx context.Context, appId string, tenantId string, data any, meta any) (*dbevent.OutboxEvent, error) {
	var metaMap map[string]any
	var err error
	if meta == nil {
		metaMap = nil
	} else if val, ok := meta.(map[string]any); ok {
		metaMap = val
	} else {
		metaMap, err = maputils.NewMap(meta)
		if err != nil {
			return nil, err
		}
	}

	_, topic, _ := reflectutils.GetTypeDetails(data)
	topic = stringutils.MidlineString(topic)
	event := &dbevent.OutboxEvent{
		Id:          idutils.NewId(),
		AppId:       appId,
		TenantId:    tenantId,
		Topic:       topic,
		Data:        data,
		Meta:        metaMap,
		CreatedTime: times.NewTime(),
	}
	return event, err
}
