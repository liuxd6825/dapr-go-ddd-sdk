package ddd

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
)

type NewEventFunc func() interface{}

var _eventTypeRegistry = newEventTypeRegistry()

type RegisterEventTypeOptions struct {
	marshaler JsonMarshaler
}

type RegisterOption func(*RegisterEventTypeOptions)

func RegisterOptionMarshaler(marshaler JsonMarshaler) RegisterOption {
	return func(options *RegisterEventTypeOptions) {
		options.marshaler = marshaler
	}
}

func RegisterEventType(eventType string, eventRevision string, newFunc NewEventFunc, options ...RegisterOption) error {
	if err := assert.NotEmpty(eventType, assert.WidthOptionsError("ddd.RegisterEventType() eventType is nil")); err != nil {
		return err
	}
	if err := assert.NotEmpty(eventRevision, assert.WidthOptionsError("ddd.RegisterEventType() eventType is nil")); err != nil {
		return err
	}
	if err := assert.NotNil(newFunc, assert.WidthOptionsError("ddd.RegisterEventType() newFunc is nil")); err != nil {
		return err
	}
	return _eventTypeRegistry.add(eventType, eventRevision, newFunc, options...)
}

func NewDomainEvent(record *EventRecord) (interface{}, error) {
	if eventTypes, ok := _eventTypeRegistry.typeMap[record.EventType]; ok {
		if item, ok := eventTypes.revisionMap[record.EventRevision]; ok {
			event := item.newFunc()
			var err error
			if item.marshaler != nil {
				err = item.marshaler(record, event)
			} else {
				err = record.Marshal(event)
			}
			if err != nil {
				_, _ = applog.Error("", "ddd", "NewDomainEvent", err.Error())
				return nil, err
			}
			return event, nil
		}
	}
	err := errors.New(fmt.Sprintf("没有注册的事件类型 %s %s", record.EventType, record.EventRevision))
	_, _ = applog.Error("", "ddd", "NewDomainEvent", err.Error())
	return nil, err
}

func getRegistryItem(eventType, eventRevision string) (*registryItem, error) {
	if eventTypes, ok := _eventTypeRegistry.typeMap[eventType]; ok {
		if item, ok := eventTypes.revisionMap[eventRevision]; ok {
			return item, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("没有注册的事件类型 %s %s", eventType, eventRevision))
}

type JsonMarshaler func(record *EventRecord, event interface{}) error

type registryItem struct {
	eventType      string
	revision       string
	newFunc        NewEventFunc
	marshaler      JsonMarshaler
	eventPrototype interface{}
}

func newRegistryItem(eventType, revision string, eventFunc NewEventFunc, eventPrototype interface{}) *registryItem {
	return &registryItem{
		eventType:      eventType,
		revision:       revision,
		newFunc:        eventFunc,
		eventPrototype: eventPrototype,
	}
}

// 事件类型注册表
type eventTypeRegistry struct {
	typeMap map[string]*eventTypes
}

//
//  newEventTypeRegistry
//  @Description: 新建事件类型注册表
//  @return *eventTypeRegistry
//
func newEventTypeRegistry() *eventTypeRegistry {
	return &eventTypeRegistry{
		typeMap: make(map[string]*eventTypes),
	}
}

//
// add
// @Description: 添加事件类型
// @receiver r
// @param eventType 事件类型
// @param revision 事件版本号
// @param newFunc 事件方法
// @param options 选项
// @return error 错误
//
func (r *eventTypeRegistry) add(eventType string, revision string, newFunc NewEventFunc, options ...RegisterOption) error {
	opts := &RegisterEventTypeOptions{}
	for _, item := range options {
		item(opts)
	}
	eventTypes, ok := r.typeMap[eventType]
	if !ok {
		ts := newEventType(eventType)
		ts.revisionMap[revision] = newRegistryItem(eventType, revision, newFunc, nil)
		r.typeMap[ts.eventType] = ts
	} else {
		_, ok := eventTypes.revisionMap[eventType]
		if !ok {
			eventTypes.revisionMap[revision] = newRegistryItem(eventType, revision, newFunc, nil)
		} else {
			return errors.New(fmt.Sprintf("%s.%s已经存存", eventType, revision))
		}
	}
	return nil
}

// 事件类
type eventTypes struct {
	eventType   string
	revisionMap map[string]*registryItem
}

func newEventType(eventType string) *eventTypes {
	return &eventTypes{
		eventType:   eventType,
		revisionMap: make(map[string]*registryItem),
	}
}

// 事件版本
type eventRevisions struct {
	domainTypes map[string]*registryItem
}

func newRevisions() *eventRevisions {
	return &eventRevisions{
		domainTypes: make(map[string]*registryItem),
	}
}
