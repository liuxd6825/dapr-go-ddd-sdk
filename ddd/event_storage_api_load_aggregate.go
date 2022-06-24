package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
)

type LoadAggregateOptions struct {
	eventStorageKey string
}
type LoadAggregateOption func(*LoadAggregateOptions)

func LoadAggregateKey(eventStorageKey string) LoadAggregateOption {
	return func(options *LoadAggregateOptions) {
		options.eventStorageKey = eventStorageKey
	}
}

//
// LoadAggregate
// @Description: 加载聚合根
// @param ctx 上下文
// @param tenantId 租户id
// @param aggregateId 聚合根id
// @param aggregate 聚合根对象
// @param opts 可选参数
// @return agg    聚合根对象
// @return isFound 是否找到
// @return err 错误
//
func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate, opts ...LoadAggregateOption) (agg Aggregate, isFound bool, err error) {
	logInfo := &applog.LogInfo{
		TenantId:  tenantId,
		ClassName: "ddd",
		FuncName:  "LoadAggregate",
		Message:   fmt.Sprintf("aggregateId=%s", aggregateId),
		Level:     applog.INFO,
	}

	_ = applog.DoAppLog(ctx, logInfo, func() (interface{}, error) {
		options := &LoadAggregateOptions{
			eventStorageKey: "",
		}
		for _, item := range opts {
			item(options)
		}
		eventStorage, e := GetEventStorage(options.eventStorageKey)
		if e != nil {
			agg, isFound, err = nil, false, e
			return agg, err
		}
		agg, isFound, err = eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
		return agg, err
	})
	return
}
