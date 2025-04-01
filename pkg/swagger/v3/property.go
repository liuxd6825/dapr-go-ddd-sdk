package swagger3

import (
	"fmt"
	"strings"
)

type Type = string

const (
	TypeNull     Type = "null"
	TypeObject        = "object"
	TypeString        = "string"
	TypeInteger       = "integer"
	TypeNumber        = "number"
	TypeArray         = "array"
	TypeBoolean       = "boolean"
	TypeDate          = "date"
	TypeDatetime      = "datetime"
	TypeBytes         = "bytes"
)

type Format = string

const (
	FormatNull                Format = ""
	FormatDateTime                   = "date-time" // 2018-11-13 20:20:39
	FormatTime                       = "time"      // 20:20:39 00:00
	FormatDate                       = "date"      // 2018-11-13
	FormatDuration                   = "duration"
	FormatEmail                      = "email"
	FormatIdnEmail                   = "idn-email"
	FormatHostname                   = "hostname"
	FormatIdnHostname                = "idn-hostname"
	FormatIpv4                       = "ipv4"
	FormatIpv6                       = "ipv6"
	FormatUuid                       = "uuid"
	FormatUri                        = "uri"
	FormatUriReference               = "uri-reference"
	FormatIri                        = "iri"
	FormatIriReference               = "iri-reference"
	FormatUriTemplate                = "uri-template"
	FormatJsonPointer                = "json-pointer"
	FormatRelativeJsonPointer        = "relative-json-pointer"
	FormatRegex                      = "regex"
)

type Property struct {
	Name        string     `json:"-"`
	Type        any        `json:"type,omitempty"`
	Title       string     `json:"title,omitempty"`
	Format      Format     `json:"format,omitempty"`
	Pattern     string     `json:"pattern,omitempty"`
	Ref         string     `json:"$ref,omitempty"`
	Description string     `json:"description,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
	Required    []string   `json:"required,omitempty"`
	Items       *Property  `json:"items,omitempty"` // nil or []*Schema or *Schema
	Minimum     *int       `json:"minimum,omitempty"`
	Maximum     *int       `json:"maximum,omitempty"`
	MinLength   *int       `json:"minLength,omitempty"`
	MaxLength   *int       `json:"maxLength,omitempty"`
	ReadOnly    bool       `json:"readOnly,omitempty"`
	WriteOnly   bool       `json:"writeOnly,omitempty"`
	Examples    []any      `json:"examples,omitempty"`
	Deprecated  bool       `json:"deprecated,omitempty"`
	Field       *Field     `json:"field,omitempty"`
	types       map[string]string
}

// Field 数据库字段
type Field struct {
	PrimaryKey   bool   `json:"primaryKey,omitempty"`
	Unique       bool   `json:"unique,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	NotNull      bool   `json:"notNull,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Size         *int   `json:"size,omitempty"`
}

type Properties map[string]*Property

func (p *Properties) Init(schema ISchema) {
	for key, value := range *p {
		value.Name = key
		value.Init(schema)
	}
}

func (p *Property) GetItems() *Property {
	return p.Items
}

func (p *Property) GetTitle() string {
	if p.Title == "" {
		return p.Name
	}
	return p.Title
}

func (p *Property) Init(schema ISchema) {
	if p.Properties != nil {
		p.Properties.Init(p)
	}
	if p.Items != nil {
		p.Items.Init(p)
	}
}

func (p *Property) GetProperties() Properties {
	return p.Properties
}

func (p *Property) SetProperties(val Properties) {
	p.Properties = val
	if val != nil {
		val.Init(p)
	}
}

func (p *Property) GetRequired() []string {
	return p.Required
}

func (p *Property) SetRequired(val []string) {
	p.Required = val
}

func (p *Property) AddRequired(val string) {
	p.Required = append(p.Required, val)
}

func (p *Property) GetType() any {
	return p.Type
}

func (p *Property) GetName() string {
	return p.Name
}

func (p *Property) GetDescription() string {
	return p.Description
}

func (p *Property) GetFormat() Format {
	return p.Format
}

func (p *Property) initTypes() {
	if p.types == nil {
		p.types = make(map[string]string)
		if v, ok := p.Type.(string); ok {
			p.types[v] = v
		} else if v, ok := p.Type.([]any); ok {
			for _, vv := range v {
				item := strings.ToLower(fmt.Sprintf("%s", vv))
				p.types[item] = item
			}
		}
	}
}

// IncludeType 是否包含的类型
func (p *Property) IncludeType(val string) bool {
	p.initTypes()
	val = strings.ToLower(val)
	_, res := p.types[val]
	return res
}

func (p *Property) GetTypes() []string {
	p.initTypes()
	var types []string
	for _, k := range p.types {
		types = append(types, k)
	}
	return types
}

func (p *Property) IsTypeNumber() bool {
	return p.IncludeType(TypeNumber)
}

func (p *Property) IsTypeBoolean() bool {
	return p.IncludeType(TypeBoolean)
}

func (p *Property) IsTypeDate() bool {
	return p.IncludeType(TypeDate)
}

func (p *Property) IsTypeDatetime() bool {
	return p.IncludeType(TypeDatetime)
}

func (p *Property) IsTypeNull() bool {
	return p.IncludeType(TypeNull)
}

func (p *Property) IsTypeObject() bool {
	return p.IncludeType(TypeObject)
}

func (p *Property) IsTypeString() bool {
	return p.IncludeType(TypeString)
}

func (p *Property) IsTypeArray() bool {
	return p.IncludeType(TypeArray)
}

func (p *Property) IsTypeInteger() bool {
	return p.IncludeType(TypeInteger)
}
