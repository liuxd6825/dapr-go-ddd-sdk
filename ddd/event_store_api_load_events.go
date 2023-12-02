package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
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
func LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest, eventStorageKey string) (resp *daprclient.LoadEventsResponse, respErr error) {
	defer func() {
		respErr = errors.GetRecoverError(respErr, recover())
	}()

	logInfo := &applog.LogInfo{
		TenantId:  req.TenantId,
		ClassName: "ddd",
		FuncName:  "LoadAggregate",
		Message:   fmt.Sprintf("%v", req),
		Level:     logs.InfoLevel,
	}

	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		eventStorage, e := GetEventStore(eventStorageKey)
		if e != nil {
			resp, respErr = nil, e
			return nil, respErr
		}
		resp, respErr = eventStorage.LoadEvent(ctx, req)
		return resp, respErr
	})
	return
}
