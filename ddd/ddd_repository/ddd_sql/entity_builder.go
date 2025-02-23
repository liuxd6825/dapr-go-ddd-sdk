package ddd_sql

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type MapEntityBuilder[T map[string]any] struct {
}

func NewMapEntityBuilder() ddd.EntityBuilder[map[string]any] {
	return &MapEntityBuilder[map[string]any]{}
}

func (m *MapEntityBuilder[T]) NewEntity() T {
	data := map[string]any{}
	return T(data)
}

func (m *MapEntityBuilder[T]) NewEntityList() []T {
	list := []T{}
	return list
}

func (m *MapEntityBuilder[T]) GetTenantId(entity T) string {
	v, ok := entity["tenantId"]
	if ok {
		return v.(string)
	}
	return ""
}

func (m *MapEntityBuilder[T]) SetTenantId(entity T, tenantId string) {
	entity["tenantId"] = tenantId
}

func (m *MapEntityBuilder[T]) GetId(entity T) string {
	v, ok := entity["id"]
	if ok {
		return v.(string)
	}
	return ""
}

func (m *MapEntityBuilder[T]) SetId(entity T, id string) {
	entity["id"] = id
}
