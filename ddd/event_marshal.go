package ddd

type EventMarshal interface {
	Marshal(record *EventRecord) error
}
