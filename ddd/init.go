package ddd

import (
	"errors"
)

var eventStorages = map[string]EventStorage{"": NewEmptyEventStorage()}
var subscribes = make([]Subscribe, 0)
var subscribeHandlers = make([]SubscribeHandler, 0)

func Start() error {
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

func GetSubscribes() []Subscribe {
	return subscribes
}

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

func RegisterEventStorage(key string, es EventStorage) {
	eventStorages[key] = es
}

func RegisterSubscribeHandler(subHandler SubscribeHandler) error {
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
