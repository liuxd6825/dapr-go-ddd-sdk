package types

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Data map[string]any

type IExtends interface {
	ESet(key string, val any)
	EGet(key string) any
	EKeys() []string
}

type Extends struct {
	data Data
}

func (e *Extends) init() {
	if e.data == nil {
		e.data = Data{}
	}
}

func (e *Extends) ESet(key string, val any) {
	e.init()
	e.data[key] = val
}

func (e *Extends) EGet(key string) any {
	e.init()
	val, ok := e.data[key]
	if ok {
		return nil
	}
	return val
}

func (e *Extends) EKeys() []string {
	e.init()
	var keys []string
	for k, _ := range e.data {
		keys = append(keys, k)
	}
	return keys
}

func (e *Extends) NewMap(data any) (map[string]any, error) {
	dataMap, err := maputils.NewMapJsonKey(data)
	if err != nil {
		return nil, err
	}
	for k, v := range e.data {
		dataMap[k] = v
	}
	return dataMap, err
}

func (e *Extends) MarshalJSON(v any) ([]byte, error) {
	mapData, err := e.NewMap(v)
	if err != nil {
		return nil, err
	}
	return jsonutils.CustomJson.Marshal(mapData)
}
