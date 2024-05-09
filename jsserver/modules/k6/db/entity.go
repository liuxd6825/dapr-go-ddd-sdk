package db

import (
	"encoding/json"
	"fmt"
)

type Entity struct {
	Values map[string]any `json:"-"`
}

const (
	TenantIdKey = "tenantId"
	IdKey       = "id"
)

func NewEntity() *Entity {
	return &Entity{Values: make(map[string]any)}
}

func (e *Entity) GetTenantId() string {
	return e.GetString(TenantIdKey)
}

func (e *Entity) SetTenantId(val string) {
	e.Set(TenantIdKey, val)
}

func (e *Entity) GetId() string {
	return e.GetString(IdKey)
}

func (e *Entity) SetId(val string) {
	e.Set(IdKey, val)
}

func (e *Entity) Get(key string) (any, bool) {
	val, ok := e.Values[key]
	return val, ok
}

func (e *Entity) Set(key string, val any) {
	e.Values[key] = val
}

func (e *Entity) GetString(key string) string {
	val, ok := e.Get(key)
	if ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%s", val)
	}
	return ""
}

func (e *Entity) GetKeys() []string {
	j := 0
	keys := make([]string, len(e.Values))
	for k := range e.Values {
		keys[j] = k
		j++
	}
	return keys
}

func (e *Entity) GetMapValues() map[string]any {
	return e.Values
}

func (e *Entity) SetMapValues(vals map[string]any) {
	e.Values = vals
}

func (e *Entity) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Values)
}
