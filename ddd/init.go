package ddd

import (
	"errors"
)

// SubscribeController type SubscribeHandler = func(Subscribe)
type SubscribeController interface {
	GetSubscribes() (*[]Subscribe, error)
	RegisterSubscribe(Subscribe) error
}

var eventStorages = map[string]EventStorage{"": NewEmptyEventStorage()}
var subscribes = make([]Subscribe, 0)
var subscribeControllers = make([]SubscribeController, 0)

func RegisterEventStorage(key string, es EventStorage) {
	eventStorages[key] = es
}

func RegisterSubscribe(ctl SubscribeController) error {
	subscribeControllers = append(subscribeControllers, ctl)
	items, err := ctl.GetSubscribes()
	if err != nil {
		return err
	}
	for _, s := range *items {
		subscribes = append(subscribes, s)
	}
	return nil
}

func Start() error {
	for _, ctl := range subscribeControllers {
		items, err := ctl.GetSubscribes()
		if err != nil {
			return err
		}
		for _, subscribe := range *items {
			if err := ctl.RegisterSubscribe(subscribe); err != nil {
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
