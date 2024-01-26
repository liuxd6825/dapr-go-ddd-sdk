package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type LoadAggregateOptions struct {
	eventStorageKey string
}

func NewLoadAggregateOptions() *LoadAggregateOptions {
	return &LoadAggregateOptions{}
}

func (o *LoadAggregateOptions) SetEventStorageKey(eventStorageKey string) *LoadAggregateOptions {
	o.eventStorageKey = eventStorageKey
	return o
}

func (o *LoadAggregateOptions) Merge(opts ...*LoadAggregateOptions) *LoadAggregateOptions {
	for _, item := range opts {
		o.eventStorageKey = item.eventStorageKey
	}
	return o
}

// LoadAggregate
// @Description: 加载聚合根
// @param ctx 上下文
// @param tenantId 租户id
// @param aggregateId 聚合根id
// @param agg 聚合根对象
// @param opts 可选参数
// @return agg    聚合根对象
// @return isFound 是否找到
// @return err 错误
func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any, opts ...*LoadAggregateOptions) (agg Aggregate, isFound bool, err error) {
	options := NewLoadAggregateOptions().Merge(opts...)
	err = logs.DebugStart(ctx, "", nil, func() error {
		eventStorage, e := GetEventStore(options.eventStorageKey)
		if e != nil {
			agg, isFound, err = nil, false, e
		}
		agg, isFound, err = eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
		return err
	})
	return
}
