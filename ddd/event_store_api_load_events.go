package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

// LoadEvents
// @Description: 获取领域事件
// @param ctx 上下文
// @param req 传入参数
// @param EventStoreKey 事件存储器key
// @return resp 响应体
// @return err 错误
func LoadEvents(ctx context.Context, req *dapr.LoadEventsRequest, eventStorageKey string) (resp *dapr.LoadEventsResponse, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	fields := logs.Fields{
		"tenantId":  req.TenantId,
		"className": "ddd",
		"funcName":  "LoadAggregate",
		"message":   fmt.Sprintf("%v", req),
		"level":     logs.InfoLevel,
	}

	err = logs.DebugStart(ctx, req.TenantId, fields, func() error {
		eventStorage, e := GetEventStore(eventStorageKey)
		if e != nil {
			resp, err = nil, e
			return err
		}
		resp, err = eventStorage.LoadEvent(ctx, req)
		return err
	})

	return
}
