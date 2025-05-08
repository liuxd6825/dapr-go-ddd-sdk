package types

import cmap "github.com/orcaman/concurrent-map"

type CMap[T any] struct {
	data cmap.ConcurrentMap
}

func NewCMap[T any]() *CMap[T] {
	return &CMap[T]{
		data: cmap.New(),
	}
}

func (m *CMap[T]) Set(key string, value T) {
	m.data.Set(key, value)
}

func (m *CMap[T]) Get(key string) (T, bool) {
	if v, ok := m.data.Get(key); ok {
		return v.(T), ok
	}
	var null T
	return null, false
}

func (m *CMap[T]) Remove(key string) {
	m.data.Remove(key)
}

func (m *CMap[T]) IsEmpty() bool {
	return m.data.IsEmpty()
}

func (m *CMap[T]) Count() int {
	return m.data.Count()
}

func (m *CMap[T]) Has(key string) bool {
	return m.data.Has(key)
}

func (m *CMap[T]) Items() map[string]any {
	data := m.data.Items()
	return data
}

func (m *CMap[T]) Keys() []string {
	return m.data.Keys()
}
