package ddd

type SubscribeHandler = func(Subscribe)
type SubscribeHandlers interface {
	GetSubscribes() []Subscribe
	RegisterSubscribe(Subscribe)
}

var eventStorage EventStorage = NewEmptyEventStorage()
var subscribes = make([]Subscribe, 0)
var subscribeController = make([]SubscribeHandlers, 0)

func RegisterEventStorage(es EventStorage) {
	eventStorage = es
}

func RegisterSubscribe(ctl SubscribeHandlers) {
	subscribeController = append(subscribeController, ctl)
	for _, s := range ctl.GetSubscribes() {
		subscribes = append(subscribes, s)
	}
}

func Start() {
	for _, ctl := range subscribeController {
		for _, subscribe := range ctl.GetSubscribes() {
			ctl.RegisterSubscribe(subscribe)
		}
	}
}

func GetSubscribes() []Subscribe {
	return subscribes
}

func GetEventStorage() EventStorage {
	return eventStorage
}
