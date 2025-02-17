package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	"io"
	"os"
	"strings"
)

type ISchema interface {
	GetName() string
	GetType() any
	GetProperties() Properties
	SetProperties(Properties)
	GetRequired() []string
	SetRequired([]string)
	AddRequired(string)
	GetItems() *Property
	Init(schema ISchema)
}

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
	AllOf      []*Ref     `json:"allOf,omitempty"`
	types      []string
	validate   *Validate
	init       bool
	loader     URLLoader
}

type Ref struct {
	Ref string `json:"$ref"`
}

type Enum struct {
}

type Item = Property

func NewSchema(fileName string, reader io.Reader) (*Schema, error) {
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
	data.FileName = fileName
	return &data, err
}

func NewSchemaWithJson(fileName string, json string) (*Schema, error) {
	var reader = strings.NewReader(json)
	return NewSchema(fileName, reader)
}

func NewSchemaFile(fileName string) (*Schema, error) {
	// 打开 JSON 文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	return NewSchema(fileName, file)
}

func (s *Schema) Init(schema ISchema) {
	if s.init {
		return
	}
	s.InitType()
	s.init = true
	if s.Properties != nil {
		s.Properties.Init(s)
	}
	if s.Items != nil {
		s.Items.Init(s)
	}
}

func (s *Schema) UseLoader(loader URLLoader) error {
	s.loader = loader
	if s.validate != nil {
		s.validate.UseLoader(s.loader)
	}
	return nil
}

func (s *Schema) GetTableName() string {
	return stringutils.AsFieldName(s.Name)
}

func (s *Schema) Compile() error {
	s.Init(s)
	if s.validate == nil {
		s.validate = NewValidate(s)

	}
	if s.loader != nil {
		s.validate.UseLoader(s.loader)
	}

	return s.validate.Compile()
}

func (s *Schema) GetJsonSchema() *jsonschema.Schema {
	return s.validate.validator
}

func (s *Schema) InitType() {
	/*
		if s.types != nil {
			return
		}
		if v, ok := s.Type.(string); ok {
			s.types = []string{v}
			if !s.NotNull {
				s.types = append(s.types, TypeNull)
			}
		} else if v, ok := s.Type.([]any); ok {
			hasNull := false
			for _, vv := range v {
				item := strings.ToLower(fmt.Sprintf("%s", vv))
				s.types = append(s.types, item)
				if vv == TypeNull {
					hasNull = true
				}
			}
			if !s.NotNull && !hasNull {
				s.types = append(s.types, TypeNull)
			}
		}*/
}

func (s *Schema) GetName() string {
	return s.Name
}

func (s *Schema) GetType() any {
	return s.Type
}

func (s *Schema) GetItems() *Property {
	return s.Items
}

func (s *Schema) GetSchema() ISchema {
	return s
}

func (s *Schema) Validate(obj any, opts ...func(*Validate) error) {
	err := s.validate.Validate(obj)
	if err != nil {
		panic(err)
	}
}

func (s *Schema) ToJson() string {
	s.Init(s)
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

/*
func (s *Schema) Convertor(obj types.Object) (map[string]any, error) {
	return s.convertor(obj, s.Properties)
}
*/

/*
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
*/

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
