package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

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
