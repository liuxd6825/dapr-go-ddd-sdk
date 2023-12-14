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
type emptyEvent struct {
	tenantId  string
	commandId string
	eventId   string
	eventType string
}

func (e *emptyEvent) GetTenantId() string {
	return e.tenantId
}

func (e *emptyEvent) GetCommandId() string {
	return e.commandId
}

func (e *emptyEvent) GetEventId() string {
	return e.eventId
}

func (e *emptyEvent) GetEventType() string {
	return e.eventType
}

func NewEmptyEvent(e Event) Event {
	if e != nil {
		return e
	}
	return &emptyEvent{}
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
