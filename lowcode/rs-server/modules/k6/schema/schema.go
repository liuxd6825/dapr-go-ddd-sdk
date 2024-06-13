package schema

import (
	schema2 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type Schema = schema2.Schema

func NewSchema(schema *Schema) *Schema {
	return schema.Init()
}

type Validate = schema2.Validate

func NewValidate(s *schema2.Schema) *schema2.Validate {
	if s.Properties != nil {
		for k, v := range s.Properties {
			v.Name = k
		}
	}
	return schema2.NewValidate(s)
}
