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

type ISchema interface {
	GetProperties() Properties
	SetProperties(Properties)

	GetRequired() []string
	SetRequired([]string)
	AddRequired(string)
}

type Schema struct {
	Id       string `json:"$id,omitempty"`
	Schema   string `json:"$schema,omitempty"`
	validate *Validate
	init     bool

	Type       string     `json:"type,omitempty"`
	Title      string     `json:"title,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Required   []string   `json:"-"`
}

type Enum struct {
}

type Types struct {
}

type Item = Property

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
		s.Properties.Init(s)
	}
	return s
}

func (s *Schema) GetSchema() ISchema {
	return s
}

func (s *Schema) Validate(obj any) error {
	s.Init()
	if s.validate == nil {
		s.validate = NewValidate(s)
	}
	return s.validate.Validate(obj)
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
		propType := prop.GetType()
		switch propType {
		case TypeObject:
			val = source.Get(key)
			obj, ok := types.AsObject(val)
			if ok {
				val, err = s.convertor(obj, prop.Properties)
			}
		case TypeBoolean:
			val, err = source.GetBool(key, prop.NotNull)
		case TypeInteger:
			val, err = source.GetInt(key, prop.NotNull)
		case TypeNumber:
			val, err = source.GetFloat(key, prop.NotNull)
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

func (s *Schema) GetProperties() Properties {
	return s.Properties
}

func (s *Schema) SetProperties(properties Properties) {
	s.Properties = properties
	if properties != nil {
		properties.Init(s)
	}
}

func (s *Schema) GetRequired() []string {
	return s.Required
}

func (s *Schema) SetRequired(v []string) {
	s.Required = v
}

func (s *Schema) AddRequired(v string) {
	s.Required = append(s.Required, v)
}
