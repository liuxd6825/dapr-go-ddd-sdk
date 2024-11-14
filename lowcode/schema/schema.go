package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"io"
	"os"
	"strings"
)

type Schema struct {
	Schema      string     `json:"$schema,omitempty"`
	Id          string     `json:"$id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type,omitempty"`
	Properties  Properties `json:"properties,omitempty"`
	Definitions Properties `json:"definitions,omitempty"`
	Required    []string   `json:"required,omitempty"`
	Default     any        `json:"default,omitempty"`
	validate    *Validate
	init        bool
}

type Type = string

const (
	TypeNull     Type = ""
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

type Enum struct {
}

type Types struct {
}

type Property struct {
	Name        string `json:"-"`
	Title       string `json:"title,omitempty"`
	Type        Type   `json:"type,omitempty"`
	Format      Format `json:"format,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
	Description string `json:"description,omitempty"`
	Ref         string `json:"$ref,omitempty"`

	// object --
	Properties            Properties `json:"properties,omitempty"`
	Required              []string   `json:"required,omitempty"`
	MaxProperties         *int
	MinProperties         *int
	PropertyNames         *Schema
	AdditionalProperties  any            // nil or bool or *Schema
	Dependencies          map[string]any // value is []string or *Schema
	DependentRequired     map[string][]string
	DependentSchemas      map[string]*Schema
	UnevaluatedProperties *Schema

	// array --
	MinItems         *int      `json:"minItems,omitempty"`
	UniqueItems      *bool     `json:"uniqueItems,omitempty"`
	ExclusiveMinimum *int      `json:"exclusiveMinimum,omitempty"`
	MaxItems         int       `json:"maxItems"`
	Contains         *Schema   `json:"contains,omitempty"`
	MinContains      *int      `json:"minContains,omitempty"`
	MaxContains      *int      `json:"maxContains,omitempty"`
	PrefixItems      []*Schema `json:"prefixItems,omitempty"`
	Items2020        *Schema   `json:"items2020,omitempty"`
	UnevaluatedItems *Schema   `json:"unevaluatedItems,omitempty"`
	Items            Items     `json:"items"`           // nil or []*Schema or *Schema
	AdditionalItems  any       `json:"additionalItems"` // nil or bool or *Schema

	// number --
	Minimum   *int `json:"minimum,omitempty"`
	Maximum   *int `json:"maximum,omitempty"`
	MinLength *int `json:"minLength,omitempty"`
	MaxLength *int `json:"maxLength,omitempty"`

	// type agnostic --
	Bool            *bool // boolean schema
	ID              string
	Anchor          string
	RecursiveRef    *Schema
	RecursiveAnchor bool
	DynamicAnchor   string // "" if not specified
	Types           *Types
	Enum            *Enum
	Const           *any
	Not             *Schema
	AllOf           []*Schema
	AnyOf           []*Schema
	OneOf           []*Schema
	If              *Schema
	Then            *Schema
	Else            *Schema

	// annotations --
	Default    *any   `json:"default,omitempty"`
	Comment    string `json:"comment,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
	WriteOnly  bool   `json:"writeOnly,omitempty"`
	Examples   []any  `json:"examples,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
}

type Items = Schema

func NewSchema(reader io.Reader) (*Schema, error) {
	var data Schema
	var err error
	// 解码 JSON 数据到结构体中
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&data); err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding JSON: %v", err))
	}
	for key, property := range data.Properties {
		property.Name = key
	}
	return &data, err
}

func NewSchemaWithJson(json string) (*Schema, error) {
	var reader = strings.NewReader(json)
	return NewSchema(reader)
}

func NewSchemaFile(pathFile string) (*Schema, error) {
	// 打开 JSON 文件
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	return NewSchema(file)
}

func (s *Schema) Init() *Schema {
	if s.init {
		return s
	}
	s.init = true
	if s.Properties != nil {
		s.Properties.Init()
	}
	if s.Definitions != nil {
		s.Definitions.Init()
	}
	return s
}

func (s *Schema) Validate(obj any) error {
	s.Init()
	if s.validate == nil {
		s.validate = NewValidate(s)
	}
	return s.validate.Validate(obj)
}

func (p *Property) GetTitle() string {
	if p.Title == "" {
		return p.Name
	}
	return p.Title
}

func (s *Schema) ToJson() string {
	s.Init()
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (s *Schema) Convertor(obj types.Object) (map[string]any, error) {
	return s.convertor(obj, s.Properties)
}

func (s *Schema) convertor(source types.Object, props Properties) (map[string]any, error) {
	s.Init()
	target := map[string]any{}
	for key, prop := range props {
		var val any
		var err error
		switch prop.Type {
		case TypeObject:
			val = source.Get(key)
			obj, ok := types.AsObject(val)
			if ok {
				val, err = s.convertor(obj, prop.Properties)
			}
		case TypeDate, TypeDatetime:
			val, err = s.convertValue(prop, source, key)
		case TypeBoolean:
			val, err = source.GetBool(key)
		case TypeInteger:
			val, err = source.GetInt(key)
		case TypeNumber:
			val, err = source.GetFloat(key)
		default:
			val, err = s.convertValue(prop, source, key)
		}
		if err != nil {
			return nil, err
		}
		target[key] = val
	}
	return target, nil
}

func (s *Schema) convertValue(prop *Property, source types.Object, key string) (val any, err error) {
	switch prop.Format {
	case FormatDateTime:
		val, err = source.GetTime(key)
	case FormatDate:
		val, err = source.GetDate(key)
	default:
		val = source.Get(key)
	}
	return val, err
}

func (p *Properties) Init() {
	for key, value := range *p {
		value.Name = key
	}
}
