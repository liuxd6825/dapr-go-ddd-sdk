package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
)

type Schema struct {
	schema.Schema
}

func NewSchema(schema *Schema) *Schema {
	if schema.Properties != nil {
		for k, v := range schema.Properties {
			v.Name = k
		}
	}
	return schema
}

type Validate = schema.Validate

func NewValidate(s *schema.Schema) *schema.Validate {
	if s.Properties != nil {
		for k, v := range s.Properties {
			v.Name = k
		}
	}
	return schema.NewValidate(s)
}
