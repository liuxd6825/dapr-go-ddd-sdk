package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
)

type Convert struct {
	Name    string         `json:"name,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}
type ConvertFunc = func(v any, sch *jsonschema.Schema, meta *MetaExtension) (any, error)

var convertMap = types.NewCMap[ConvertFunc]()

func NewConvert() *Convert {
	return &Convert{
		Options: map[string]any{},
	}
}
func (c *Convert) init(vals map[string]any) error {
	for k, v := range vals {
		if k == "name" {
			c.Name = v.(string)
		} else {
			c.Options[k] = v
		}
	}
	return nil
}

func RegisterConvert(name string, convert ConvertFunc) error {
	convertMap.Set(name, convert)
	return nil
}

func DoConvert(sch *jsonschema.Schema, data any) (any, error) {
	if sch == nil {
		return data, nil
	}
	metaSch := GetMetaExtension(sch)
	if metaSch == nil {
		return data, nil
	}
	if metaSch.Convert == nil {
		return data, nil
	}
	name := metaSch.Convert.Name
	if converter, ok := convertMap.Get(name); ok {
		return converter(data, sch, metaSch)
	}
	return data, nil
}

func init() {
	err := RegisterConvert("PagingQueryConverter", PagingQueryConverter)
	if err != nil {
		return
	}
}
