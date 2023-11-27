package ddd

import (
	"errors"
)

var _appId string
var eventStorages = map[string]EventStorage{"": NewEmptyEventStorage()}
var subscribes = make([]*Subscribe, 0)
var subscribeHandlers = make([]SubscribeHandler, 0)

func Init(appId string) {
	_appId = appId
}

func AppId() string {
	return _appId
}

//
// StartSubscribeHandlers
// @Description: 启动事件监听
// @return error
//
func StartSubscribeHandlers() error {
	for _, handler := range subscribeHandlers {
		items, err := handler.GetSubscribes()
		if err != nil {
			return err
		}
		for _, subscribe := range items {
			if err := handler.RegisterSubscribe(subscribe); err != nil {
				return err
			}
		}
	}
	return nil
}

//
// GetSubscribes
// @Description: 获取事件监听列表
// @return []Subscribe
//
func GetSubscribes() []*Subscribe {
	return subscribes
}

//
// GetEventStorage
// @Description: 获取事件存储器
// @param key 事件存储器名称
// @return EventStorage 事件存储器
// @return error
//
func GetEventStorage(key string) (EventStorage, error) {
	eventStorage, ok := eventStorages[key]
	if !ok {
		return nil, errors.New("eventStorage is nil")
	}
	if eventStorage == nil {
		return nil, errors.New("eventStorage is nil")
	}
	_, ok = eventStorage.(*emptyEventStorage)
	if ok {
		return nil, errors.New("eventStorage is EmptyEventStorage")
	}
	return eventStorage, nil
}

//
// RegisterEventStorage
// @Description:  注册领域事件存储器
// @param key 唯一名称
// @param es  事件存储器
//
func RegisterEventStorage(key string, es EventStorage) {
	eventStorages[key] = es
}

//
// RegisterQueryHandler
// @Description:  注册事件监听器
// @param subHandler
// @return error
//
func RegisterQueryHandler(subHandler SubscribeHandler) error {
	subscribeHandlers = append(subscribeHandlers, subHandler)
	items, err := subHandler.GetSubscribes()
	if err != nil {
		return err
	}
	for _, s := range items {
		subscribes = append(subscribes, s)
	}
	return nil
}
