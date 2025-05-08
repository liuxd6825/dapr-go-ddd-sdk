package schema

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
	types       []string
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

func (p *Property) IncludeType(val string) bool {
	if p.types == nil {
		p.types = []string{}
		if v, ok := p.Type.(string); ok {
			p.types = append(p.types, v)
		} else if v, ok := p.Type.([]any); ok {
			for _, vv := range v {
				item := strings.ToLower(fmt.Sprintf("%s", vv))
				p.types = append(p.types, item)
			}
		}
	}
	val = strings.ToLower(val)
	for _, v := range p.types {
		if v == val {
			return true
		}
	}
	return false
}
