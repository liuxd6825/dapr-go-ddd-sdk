package fs

type Watcher interface {
	Start(opts ...Option) error
}
type WatcherEventType int

const (
	WriteEvent WatcherEventType = iota
	RemoveEvent
	CreateEvent
	RenameEvent
	ChmodEvent
)

type Option func(rootPath, fileName string, eventType WatcherEventType) error

type WatcherFs interface {
	NewWatcher() (Watcher, error)
}
