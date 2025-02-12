package ddd_sql

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type MapEntityBuilder[T MapEntity] struct {
}

func NewMapEntityBuilder() ddd.EntityBuilder[MapEntity] {
	return &MapEntityBuilder[MapEntity]{}
}

func (m *MapEntityBuilder[T]) NewEntity() (T, error) {
	data := MapEntity{}
	return T(data), nil
}

func (m *MapEntityBuilder[T]) NewEntityList() ([]T, error) {
	list := []T{}
	return list, nil
}

func (m *MapEntityBuilder[T]) GetTenantId(entity T) string {
	v, ok := entity["tenant_id"]
	if ok {
		return v.(string)
	}
	return ""
}

func (m *MapEntityBuilder[T]) SetTenantId(entity T, tenantId string) {
	entity["tenant_id"] = tenantId
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
