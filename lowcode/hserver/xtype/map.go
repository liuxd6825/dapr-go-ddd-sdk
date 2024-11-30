package xtype

import cmap "github.com/orcaman/concurrent-map"

type Map[T any] struct {
	mapData cmap.ConcurrentMap
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		mapData: cmap.New(),
	}
}

func (m *Map[T]) Set(key string, value T) {
	m.mapData.Set(key, value)
}

func (m *Map[T]) Get(key string) (T, bool) {
	if v, ok := m.mapData.Get(key); ok {
		return v.(T), ok
	}
	var null T
	return null, false
}

func (m *Map[T]) Remove(key string) {
	m.mapData.Remove(key)
}

func (m *Map[T]) IsEmpty() bool {
	return m.mapData.IsEmpty()
}

func (m *Map[T]) Count() int {
	return m.mapData.Count()
}

func (m *Map[T]) Has(key string) bool {
	return m.mapData.Has(key)
}

func (m *Map[T]) Items() map[string]any {
	data := m.mapData.Items()
	return data
}
