package ddd

type Subscribe struct {
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic"`
	Route      string            `json:"route"`
	Metadata   map[string]string `json:"metadata"`
}

type SubscribeEventHandler interface {
	DoEvent(record EventRecord)
}

func NewSubscribeItem(pubsubName string, topic, route string, metadata map[string]string, handler interface{}) *Subscribe {
	return &Subscribe{
		PubsubName: pubsubName,
		Topic:      topic,
		Metadata:   metadata,
		Route:      route,
	}
}
