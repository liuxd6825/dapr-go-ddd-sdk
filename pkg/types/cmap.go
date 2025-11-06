package types

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	cmap "github.com/orcaman/concurrent-map"
)

type CMap[T any] struct {
	data cmap.ConcurrentMap
}

func init() {
	cmap.SHARD_COUNT = 5
}

func NewCMap[T any]() *CMap[T] {
	return &CMap[T]{
		data: cmap.New(),
	}
}

func (m *CMap[T]) MSet(v map[string]T) {
	for key, value := range v {
		key = strings.Trim(key, " ")
		m.data.Set(key, value)
	}
}

func (m *CMap[T]) Set(key string, value T) {
	key = strings.Trim(key, " ")
	m.data.Set(key, value)
}

func (m *CMap[T]) Add(key string, value T) {
	key = strings.Trim(key, " ")
	has := m.data.Has(key)
	if has {
		panic(errors.New("key already exists %s", key))
	}
	m.data.Set(key, value)
}

func (m *CMap[T]) Get(key string) (T, bool) {
	key = strings.Trim(key, " ")
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

func (m *CMap[T]) Items() map[string]T {
	data := m.data.Items()
	values := make(map[string]T, len(data))
	for key, value := range data {
		values[key] = value.(T)
	}
	return values
}

func (m *CMap[T]) Keys() []string {
	return m.data.Keys()
}

func (m *CMap[T]) Clear() {
	m.data.Clear()
}

func (m *CMap[T]) MarshalJSON() ([]byte, error) {
	return m.data.MarshalJSON()
}
