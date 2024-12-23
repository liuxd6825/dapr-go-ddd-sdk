package xschema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils/schema_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"

	"github.com/liuxd6825/jsonschema/v6"
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

type Schema struct {
	FileName   string     `json:"fileName,omitempty"`
	Ref        string     `json:"$ref,omitempty"`
	Id         string     `json:"$id,omitempty"`
	Schema     string     `json:"$schema,omitempty"`
	Type       any        `json:"type,omitempty"`
	Title      string     `json:"title,omitempty"`
	Name       string     `json:"name,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Items      *Property  `json:"items,omitempty"`
	Required   []string   `json:"required,omitempty"`
	ReadOnly   bool       `json:"readOnly,omitempty"`
	WriteOnly  bool       `json:"writeOnly,omitempty"`
	AllOf      []*Schema  `json:"allOf,omitempty"`

	jsonschema *jsonschema.Schema
}

func NewSchema() *Schema {
	return &Schema{}
}

func (sch *Schema) Init(fileName string, schemaLoader schema.URLLoader) *jsonschema.Schema {
	if sch.jsonschema != nil {
		return sch.jsonschema
	}
	scm, err := schema_utils.Compile(fileName, sch, func(c *jsonschema.Compiler) error {
		c.UseLoader(schemaLoader)
		return nil
	})
	if err != nil {
		panic(err)
	}
	sch.jsonschema = scm
	return scm
}

func (sch *Schema) GetSchema() *jsonschema.Schema {
	return sch.jsonschema
}
