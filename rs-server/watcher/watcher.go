package watcher

type Watcher interface {
	Start(opts ...Options) error
}
type EventType int

const (
	WriteEvent EventType = iota
	RemoveEvent
	CreateEvent
	RenameEvent
	ChmodEvent
)

type Options func(rootPath, fileName string, eventType EventType) error
