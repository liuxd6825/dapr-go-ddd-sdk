package schema

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
	Name          string     `json:"-"`
	Type          string     `json:"type,omitempty"`
	Title         string     `json:"title,omitempty"`
	Format        Format     `json:"format,omitempty"`
	Pattern       string     `json:"pattern,omitempty"`
	Ref           string     `json:"$ref,omitempty"`
	Description   string     `json:"description,omitempty"`
	Properties    Properties `json:"properties,omitempty"`
	MaxProperties *int       `json:"maxProperties,omitempty"`
	MinProperties *int       `json:"minProperties,omitempty"`
	Required      []string   `json:"-"`
	NotNull       bool       `json:"notNull,omitempty"`

	// array --
	MinItems         *int        `json:"minItems,omitempty"`
	UniqueItems      *bool       `json:"uniqueItems,omitempty"`
	ExclusiveMinimum *int        `json:"exclusiveMinimum,omitempty"`
	MaxItems         *int        `json:"maxItems,omitempty"`
	Contains         *Property   `json:"contains,omitempty"`
	MinContains      *int        `json:"minContains,omitempty"`
	MaxContains      *int        `json:"maxContains,omitempty"`
	PrefixItems      []*Property `json:"prefixItems,omitempty"`
	Items2020        *Property   `json:"items2020,omitempty"`
	UnevaluatedItems *Property   `json:"unevaluatedItems,omitempty"`
	Items            *Property   `json:"items,omitempty"` // nil or []*Schema or *Schema
	AdditionalItems  any         `json:"additionalItems"` // nil or bool or *Schema

	// number --
	Minimum   *int `json:"minimum,omitempty"`
	Maximum   *int `json:"maximum,omitempty"`
	MinLength *int `json:"minLength,omitempty"`
	MaxLength *int `json:"maxLength,omitempty"`

	// type agnostic --
	Bool            bool        `json:"bool,omitempty"` // boolean schema
	ID              string      `json:"id,omitempty"`
	Anchor          string      `json:"anchor,omitempty"`
	RecursiveRef    *Property   `json:"recursiveRef,omitempty"`
	RecursiveAnchor bool        `json:"recursiveAnchor,omitempty"`
	DynamicAnchor   string      `json:"dynamicAnchor"` // "" if not specified
	Types           *Types      `json:"types,omitempty"`
	Enum            *Enum       `json:"enum,omitempty"`
	Const           *any        `json:"const,omitempty"`
	Not             *Property   `json:"not,omitempty"`
	AllOf           []*Property `json:"allOf,omitempty"`
	AnyOf           []*Property `json:"anyOf,omitempty"`
	OneOf           []*Property `json:"oneOf,omitempty"`
	If              *Property   `json:"if,omitempty"`
	Then            *Property   `json:"then,omitempty"`
	Else            *Property   `json:"else,omitempty"`

	//
	Comment    string `json:"comment,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
	WriteOnly  bool   `json:"writeOnly,omitempty"`
	Examples   []any  `json:"examples,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
}

type Properties map[string]*Property

func (p *Properties) Init(schema ISchema) {
	for key, value := range *p {
		value.Name = key
		value.init(schema)
	}
}

func (p *Property) GetTitle() string {
	if p.Title == "" {
		return p.Name
	}
	return p.Title
}

func (p *Property) init(schema ISchema) {
	if p.Type == "" {
		p.Type = TypeString
	}
	switch p.Type {
	case TypeDatetime:
		p.Type = TypeString
		p.Format = FormatDateTime
	case TypeDate:
		p.Type = TypeString
		p.Format = FormatDate
	}
	if p.NotNull {
		schema.AddRequired(p.Name)
		if p.Type == TypeString && p.MinLength == nil {
			minLength := 1
			p.MinLength = &minLength
		}
	} else {
		p.OneOf = []*Property{
			&Property{Type: p.Type, Format: p.Format},
			&Property{Type: TypeNull},
		}
		p.Type = ""
	}
	if p.Properties != nil {
		p.Properties.Init(p)
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

func (p *Property) GetType() Type {
	if len(p.OneOf) > 0 {
		return p.OneOf[0].Type
	}
	return p.Type
}
