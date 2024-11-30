package xschema

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
}

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

type Properties map[string]*Property

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
