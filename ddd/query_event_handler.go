package ddd

type QueryEventHandler interface {
	OnEvent(record *EventRecord) error
}
