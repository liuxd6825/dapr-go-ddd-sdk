package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonschemautils"
	"github.com/liuxd6825/jsonschema/v6"
)

type Validate struct {
	schema    *Schema
	validator *jsonschema.Schema
	compiler  *jsonschema.Compiler
}

func NewValidate(schema *Schema) *Validate {
	return &Validate{schema: schema, compiler: jsonschemautils.NewCompiler()}
}

func (v *Validate) UseLoader(loader URLLoader) {
	v.compiler.UseLoader(loader)
}

func (v *Validate) Compile() error {
	bs, err := json.Marshal(v.schema)
	if err != nil {
		return err
	}
	data, err := jsonschema.UnmarshalJSON(bytes.NewReader(bs))
	if err != nil {
		return err
	}
	fileName := v.schema.FileName
	if fileName == "" {
		fileName = "schema.json"
	}
	if err := v.compiler.AddResource(fileName, data); err != nil {
		return err
	}
	validator, err := v.compiler.Compile(fileName)
	if err != nil {
		return err
	}
	v.validator = validator
	return nil
}

func (v *Validate) Validate(val any) error {
	if v.validator == nil {
		if err := v.Compile(); err != nil {
			var sErr *jsonschema.SchemaValidationError
			if errors.As(err, &sErr) {
				return NewSchemaError(sErr)
			}
			return err
		}
	}
	if m, ok := val.(types.Object); ok {
		val = m.AsMap()
	}
	err := v.validator.Validate(val)
	var validateError *jsonschema.ValidationError
	if errors.As(err, &validateError) {
		return NewFieldsError(validateError)
	}
	return err
}
