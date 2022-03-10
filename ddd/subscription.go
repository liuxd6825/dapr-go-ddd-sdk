package ddd

type Subscribe struct {
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic"`
	Route      string            `json:"route"`
	Metadata   map[string]string `json:"metadata"`
	Handle     interface{}
}

type SubscribeEventHandler interface {
	DoEvent(record EventRecord)
}

func NewSubscribeItem(pubsubName string, topic, route string, metadata map[string]string, handle interface{}) *Subscribe {
	return &Subscribe{
		PubsubName: pubsubName,
		Topic:      topic,
		Metadata:   metadata,
		Handle:     handle,
	}
}
