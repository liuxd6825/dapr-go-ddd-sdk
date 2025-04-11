package engine

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"strings"
)

type Property struct {
	*jsonschema.Schema
	meta *schema.MetaExtension
}

func NewProperty(sch *jsonschema.Schema) *Property {
	return &Property{
		Schema: sch,
		meta:   nil,
	}
}

func (p *Property) HasId() bool {
	if strings.HasSuffix(strings.ToLower(p.Name()), "id") {
		return true
	}
	return false
}

func (p *Property) Meta() *schema.MetaExtension {
	if p.meta == nil {
		p.meta = schema.GetMetaExtension(p.Schema)
	}
	if p.meta == nil {
		p.meta = schema.NewMetaExtension()
	}
	return p.meta
}

func (p *Property) Readonly() bool {
	return p.Meta().DBField.Readonly()
}

func (p *Property) PK() bool {
	return p.Meta().DBField.PrimaryKey
}

func (p *Property) Type() string {
	return "text"
	if p.Types.Contains(jsonschema.JsonType_BooleanType) {
		return "boolean"
	} else if p.Types.Contains(jsonschema.JsonType_NumberType) {
		return "number"
	} else if p.Types.Contains(jsonschema.JsonType_StringType) {
		return "string"
	} else if p.Types.Contains(jsonschema.JsonType_IntegerType) {
		return "integer"
	} else if p.Types.Contains(jsonschema.JsonType_DateType) {
		return "date"
	} else if p.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return "datetime"
	}
	return ""
}

func GetAllProperties(sch *jsonschema.Schema) []*Property {
	properties := make([]*Property, 0)
	for _, item := range sch.GetSortProperties() {
		prop := NewProperty(item)
		properties = append(properties, prop)
	}
	return properties
}
