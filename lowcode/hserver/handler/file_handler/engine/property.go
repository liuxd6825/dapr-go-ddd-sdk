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

func (p *Property) Hide() bool {
	if p.getColumn() != nil {
		if val, ok := p.getColumn()["hide"].(bool); ok {
			return val
		}
		if val, ok := p.getColumn()["visible"].(string); ok {
			if val == "true" {
				return true
			}
		}
	}
	return false
}

func (p *Property) Readonly() bool {
	return p.Meta().DBField.Readonly()
}

func (p *Property) PK() bool {
	return p.Meta().DBField.PrimaryKey
}

func (p *Property) ColumnValue(propName string) any {
	if p.getColumn() != nil {
		if val, ok := p.getColumn()[propName]; ok {
			return val
		}
	}
	return nil
}

func (p *Property) getColumn() schema.Column {
	if p.Meta().Column != nil {
		return *(p.Meta().Column)
	}
	return nil
}

func (p *Property) FormOptionsValue(propName string) any {
	if opts := p.FormValue("options"); opts != nil {
		optsMap := opts.(map[string]interface{})
		if val, ok := optsMap[propName]; ok {
			return val
		}
	}
	return nil
}

func (p *Property) FormValue(propName string) any {
	if p.getForm() != nil {
		if val, ok := p.getForm()[propName]; ok {
			return val
		}
	}
	return nil
}

func (p *Property) getForm() schema.Form {
	if p.Meta().Form != nil {
		return *(p.Meta().Form)
	}
	return nil
}

func (p *Property) FormCtl() string {
	if p.Types.Contains(jsonschema.JsonType_BooleanType) {
		return "ui5-select"
	} else if p.Types.Contains(jsonschema.JsonType_NumberType) {
		return "ui5-input"
	} else if p.Types.Contains(jsonschema.JsonType_StringType) {
		return "ui5-input"
	} else if p.Types.Contains(jsonschema.JsonType_IntegerType) {
		return "ui5-input"
	} else if p.Types.Contains(jsonschema.JsonType_DateType) {
		return "ui5-date-picker"
	} else if p.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return "ui5-datetime-picker"
	}
	return "ui5-input"
}

func (p *Property) Type() string {
	if p.getColumn() != nil {
		if val, ok := p.getColumn()["dataType"].(string); ok {
			return val
		}
	}
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

func (p *Property) SheetType() string {
	if p.getColumn() != nil {
		if val, ok := p.getColumn()["type"].(string); ok {
			return val
		}
	}

	if p.Types.Contains(jsonschema.JsonType_BooleanType) {
		return "checkbox"
	} else if p.Types.Contains(jsonschema.JsonType_NumberType) {
		return "numeric"
	} else if p.Types.Contains(jsonschema.JsonType_StringType) {
		return "text"
	} else if p.Types.Contains(jsonschema.JsonType_IntegerType) {
		return "numeric"
	} else if p.Types.Contains(jsonschema.JsonType_DateType) {
		return "date"
	} else if p.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return "date"
	} else {
		return "text"
	}

}

func GetAllProperties(sch *jsonschema.Schema) []*Property {
	properties := make([]*Property, 0)
	for _, item := range sch.GetSortProperties() {
		prop := NewProperty(item)
		properties = append(properties, prop)
	}
	return properties
}
