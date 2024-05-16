package ddd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
)

type EntityBuilder[T any] interface {
	NewEntity() (T, error)
	NewEntityList() ([]T, error)
	GetTenantId(entity T) string
	SetTenantId(entity T, tenantId string)
	GetId(entity T) string
	SetId(entity T, id string)
}

type MapEntityBuilder[T any] struct {
}

func NewMapEntityBuilder[T any]() EntityBuilder[T] {
	return &MapEntityBuilder[T]{}
}

func (m *MapEntityBuilder[T]) NewEntity() (T, error) {
	mv := NewMapEntity()
	var a any = mv
	if t, ok := a.(T); ok {
		return t, nil
	}
	var null T
	return null, fmt.Errorf("cannot convert map to entity")
}

func (m *MapEntityBuilder[T]) NewEntityList() ([]T, error) {
	return make([]T, 0), nil
}

func (m *MapEntityBuilder[T]) GetTenantId(entity T) string {
	return get(entity, "tenantId")
}

func (m *MapEntityBuilder[T]) SetTenantId(entity T, id string) {
	m.Set(entity, "tenantId", id)
}

func (m *MapEntityBuilder[T]) GetId(entity T) string {
	return get(entity, "id")
}

func (m *MapEntityBuilder[T]) SetId(entity T, id string) {
	m.Set(entity, "id", id)
}

func get(entity any, key string) string {
	var a any = entity
	if mv, ok := a.(MapEntity); ok {
		if val, ok := mv[key]; ok {
			if v, ok := val.(string); ok {
				return v
			} else {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	return ""
}

func (m *MapEntityBuilder[T]) Set(entity T, key string, value string) {
	var a any = entity
	if mv, ok := a.(MapEntity); ok {
		mv[key] = value
	}
}

func (m *MapEntityBuilder[T]) as(entity T) (map[string]any, bool) {
	var a any = entity
	if v, ok := a.(MapEntity); ok {
		return v, true
	}
	return nil, false
}

type StructEntityBuilder[T any] struct {
}

func NewStructEntityBuilder[T any]() EntityBuilder[T] {
	return &StructEntityBuilder[T]{}
}

func (s *StructEntityBuilder[T]) NewEntity() (T, error) {
	return reflectutils.NewStruct[T]()
}

func (s *StructEntityBuilder[T]) NewEntityList() ([]T, error) {
	return reflectutils.NewSlice[[]T]()
}

func (s *StructEntityBuilder[T]) GetTenantId(entity T) string {
	if e, ok := s.as(entity); ok {
		return e.GetTenantId()
	}
	return ""
}

func (s *StructEntityBuilder[T]) SetTenantId(entity T, tenantId string) {
	if e, ok := s.as(entity); ok {
		e.SetTenantId(tenantId)
	}
}

func (s *StructEntityBuilder[T]) GetId(entity T) string {
	if e, ok := s.as(entity); ok {
		return e.GetId()
	}
	return ""
}

func (s *StructEntityBuilder[T]) SetId(entity T, id string) {
	if e, ok := s.as(entity); ok {
		e.SetId(id)
	}
}

func (s *StructEntityBuilder[T]) as(entity T) (Entity, bool) {
	var a any = entity
	if e, ok := a.(Entity); ok {
		return e, true
	}
	return nil, false
}
