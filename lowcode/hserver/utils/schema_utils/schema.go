package schema_utils

import (
	"bytes"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema/formats"
	"github.com/liuxd6825/jsonschema/v6"
)

func Compile(fileName string, schemaConfig any, opts ...func(comilper *jsonschema.Compiler) error) (schemaCompiler *jsonschema.Schema, err error) {
	compiler := newCompiler()
	compiler.AssertFormat()
	compiler.AssertContent()

	for _, opt := range opts {
		if opt != nil {
			if err = opt(compiler); err != nil {
				return nil, err
			}
		}
	}

	bs, err := json.Marshal(schemaConfig)
	if err != nil {
		return nil, err
	}
	data, err := jsonschema.UnmarshalJSON(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	if err := compiler.AddResource(fileName, data); err != nil {
		return nil, err
	}
	newSchema, err := compiler.Compile(fileName)
	if err != nil {
		return nil, err
	}
	return newSchema, nil
}

func newCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	compiler.AssertContent()
	compiler.RegisterFormat(formats.DateTimeFormat)
	compiler.RegisterFormat(formats.DateFormat)
	return compiler
}
