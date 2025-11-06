package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
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
func (c *Convert) init(ctx *jsonschema.CompilerContext, vals map[string]any) error {
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
	if err := Validate(sch, data); err != nil {
		return nil, err
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

func Validate(sch *jsonschema.Schema, data any) error {
	err := sch.Validate(data)
	if e, ok := err.(*jsonschema.ValidationError); ok {
		err = NewFieldsError(sch, e)
	} else if e, ok := err.(*jsonschema.SchemaValidationError); ok {
		err = NewSchemaError(e)
	}
	return err
}

func init() {
	err := RegisterConvert("PagingQueryConverter", PagingQueryConverter)
	if err != nil {
		return
	}
}
