package ddd

import (
	"errors"
	"fmt"
)

type NewEventFunc func() interface{}

var _registry = newEventRegistry()

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
	return _registry.add(eventType, eventRevision, newFunc, options...)
}

func NewDomainEvent(record *EventRecord) (interface{}, error) {
	if eventTypes, ok := _registry.typeMap[record.EventType]; ok {
		if item, ok := eventTypes.revisionMap[record.EventRevision]; ok {
			event := item.newFunc()
			var err error
			if item.marshaler != nil {
				err = item.marshaler(record, event)
			} else {
				err = record.Marshal(event)
			}
			if err != nil {
				return nil, err
			}
			return event, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("没有注册的事件类型 %s %s", record.EventType, record.EventRevision))
}

type JsonMarshaler func(record *EventRecord, event interface{}) error

type registryItem struct {
	eventType string
	revision  string
	newFunc   NewEventFunc
	marshaler JsonMarshaler
}

func newRegistryItem(eventType, revision string, eventFunc NewEventFunc) *registryItem {
	return &registryItem{
		eventType: eventType,
		revision:  revision,
		newFunc:   eventFunc,
	}
}

// 事件注册表
type eventRegistry struct {
	typeMap map[string]*eventTypes
}

func newEventRegistry() *eventRegistry {
	return &eventRegistry{
		typeMap: make(map[string]*eventTypes),
	}
}

func (r *eventRegistry) add(eventType string, revision string, newFunc NewEventFunc, options ...RegisterOption) error {
	opts := &RegisterEventTypeOptions{}
	for _, item := range options {
		item(opts)
	}
	eventTypes, ok := r.typeMap[eventType]
	if !ok {
		ts := newEventType(eventType)
		ts.revisionMap[revision] = newRegistryItem(eventType, revision, newFunc)
		r.typeMap[ts.eventType] = ts
	} else {
		_, ok := eventTypes.revisionMap[eventType]
		if !ok {
			eventTypes.revisionMap[revision] = newRegistryItem(eventType, revision, newFunc)
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
