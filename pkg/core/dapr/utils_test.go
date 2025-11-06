package dapr

import "testing"

type Event struct {
	EventType string `json:"eventType"`
	EventId   string `json:"eventId"`
	EventTime string `json:"eventTime"`
	Data      any    `json:"data"`
}

func TestGetTopic(t *testing.T) {
	topic1 := GetTopic(Event{})
	t.Log(topic1)

	topic2 := GetTopic(Event{})
	t.Log(topic2)
}
