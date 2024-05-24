package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema/formats"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Validate struct {
	schema   *Schema
	validate *jsonschema.Schema
}

func NewValidate(schema *Schema) *Validate {
	return &Validate{schema: schema}
}

func (v *Validate) Compiler() error {
	bs, err := json.Marshal(v.schema)
	if err != nil {
		return err
	}
	jsschema, err := jsonschema.UnmarshalJSON(bytes.NewReader(bs))
	if err != nil {
		return err
	}
	c := jsonschema.NewCompiler()
	c.AssertFormat()
	c.AssertContent()
	c.RegisterFormat(formats.DateTimeFormat)
	c.RegisterFormat(formats.DateFormat)
	if err := c.AddResource("schema.json", jsschema); err != nil {
		return err
	}
	validate, err := c.Compile("schema.json")
	if err != nil {
		return err
	}
	v.validate = validate
	return nil
}

func (v *Validate) Validate(val any) error {
	if v.validate == nil {
		if err := v.Compiler(); err != nil {
			var sErr *jsonschema.SchemaValidationError
			if errors.As(err, &sErr) {
				return NewSchemaError(sErr)
			}
			return err
		}
	}
	if m, ok := val.(common.Object); ok {
		val = m.AsMap()
	}
	err := v.validate.Validate(val)
	var validateError *jsonschema.ValidationError
	if errors.As(err, &validateError) {
		return NewFieldsError(validateError)
	}
	return err
}
