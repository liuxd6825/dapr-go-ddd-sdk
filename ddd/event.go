package ddd

import (
	"encoding/json"
	"time"
)

type DomainEvent interface {
	GetTenantId() string       // 租户Id
	GetCommandId() string      // 命令Id
	GetEventId() string        // 事件Id
	GetEventType() string      // 事件类型
	GetEventVersion() string   // 事件版本号
	GetAggregateId() string    // 聚合根Id
	GetCreatedTime() time.Time // 创建时间
	GetData() interface{}      // 事件数据
}

type IsSourcing interface {
	GetIsSourcing() bool // 是否需要事件溯源
}

type Event interface {
	GetTenantId() string  // 租户Id
	GetCommandId() string // 命令Id
	GetEventId() string   // 事件Id
	GetEventType() string // 事件类型
}

func DoEvent(eventData map[string]interface{}, event interface{}) *DoEventResult {
	bs, err := json.Marshal(eventData)
	err = json.Unmarshal(bs, event)
	result := &DoEventResult{err: err, event: event}
	return result
}

type DoEventResult struct {
	err   error
	event interface{}
}

func (r *DoEventResult) OnSuccess(fun func(event interface{}) error) *DoEventResult {
	err := fun(r.event)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *DoEventResult) OnError(fun func(err error)) *DoEventResult {
	fun(r.err)
	return r
}

func (r *DoEventResult) Error() error {
	return r.err
}
