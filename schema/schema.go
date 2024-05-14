package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
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
	validate    *Validate
}

type Type = string

const (
	TypeNull    Type = ""
	TypeObject       = "object"
	TypeString       = "string"
	TypeInteger      = "integer"
	TypeNumber       = "number"
	TypeArray        = "array"
	TypeBoolean      = "boolean"
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
	Name             string     `json:"-"`
	Title            string     `json:"title,omitempty"`
	Type             Type       `json:"type,omitempty"`
	Format           Format     `json:"format,omitempty"`
	Pattern          string     `json:"pattern,omitempty"`
	Properties       Properties `json:"properties,omitempty"`
	Required         []string   `json:"required,omitempty"`
	Description      string     `json:"description,omitempty"`
	Ref              string     `json:"$ref,omitempty"`
	Items            *Items     `json:"items,omitempty"`
	MinItems         *int       `json:"minItems,omitempty"`
	UniqueItems      *bool      `json:"uniqueItems,omitempty"`
	ExclusiveMinimum *int       `json:"exclusiveMinimum,omitempty"`
	Minimum          *int       `json:"minimum,omitempty"`
	Maximum          *int       `json:"maximum,omitempty"`
}

type Items struct {
	Type string `json:"type"`
}

func NewSchemaString(json string) (*Schema, error) {
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

func (s *Schema) Validate(obj any) error {
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
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (s *Schema) Convertor(source common.Object) (common.Object, error) {
	return s.convertor(source, s.Properties)
}

func (s *Schema) convertor(source common.Object, props Properties) (common.Object, error) {
	target := common.NewObject()
	for key, prop := range props {
		var val any
		var err error
		switch prop.Type {
		case TypeObject:
			val = source.Get(key)
			obj, ok := common.AsObject(val)
			if ok {
				val, err = s.convertor(obj, prop.Properties)
			}
		case TypeBoolean:
			val, err = source.GetBool(key)
		case TypeInteger:
			val, err = source.GetInt(key)
		case TypeNumber:
			val, err = source.GetFloat(key)
		case TypeString:
			switch prop.Format {
			case FormatDateTime:
				val, err = source.GetDateTime(key)
			case FormatDate:
				val, err = source.GetDateTime(key)
			default:
				val = source.Get(key)
				err = nil
			}
		default:
			val = source.Get(key)
			err = nil
		}
		if err != nil {
			return nil, err
		}
		_ = target.Set(key, val)

	}
	return target, nil
}
