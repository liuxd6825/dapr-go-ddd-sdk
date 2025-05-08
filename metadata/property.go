package metadata

import "reflect"
import "strings"

type Property interface {
	Name() string
	TypeName() string
	Description() string
	Json() string
	IsArray() bool
	Kind() string
}

type property struct {
	name        string              `json:"name"`
	json        string              `json:"json"`
	typeName    string              `json:"typeName,omitempty"`
	fieldInfo   reflect.StructField `json:"-"`
	description string              `json:"title"`
	isArray     bool                `json:"isArray"`
	kind        string              `json:"kind"`
}

func NewProperty() Property {
	return &property{}
}

func (p *property) Init(fieldInfo reflect.StructField) Property {
	p.fieldInfo = fieldInfo
	p.name = fieldInfo.Name
	p.typeName = p.getTypeName(fieldInfo)
	p.isArray = p.getIsArray()
	p.json = p.getJson(fieldInfo)
	p.description = p.getDesc(fieldInfo)
	p.kind = fieldInfo.Type.Kind().String()
	return p
}

func (p *property) getIsArray() bool {
	switch p.fieldInfo.Type.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	}
	return false
}

func (p *property) getDesc(fieldInfo reflect.StructField) string {
	if desc := fieldInfo.Tag.Get(DescTagName); desc != "" {
		return desc
	} else if val := fieldInfo.Tag.Get(DescriptionTagName); val != "" {
		return val
	}
	return ""
}

func (p *property) getJson(fieldInfo reflect.StructField) string {
	if json := fieldInfo.Tag.Get(JsonTagName); json != "" {
		json = strings.Replace(json, ",omitempty", "", 100)
		return json
	}
	return fieldInfo.Name
}

func (p *property) getTypeName(fieldInfo reflect.StructField) string {
	res := fieldInfo.Type.Name()
	if fieldInfo.Type.Kind() == reflect.Ptr {
		res = fieldInfo.Type.Elem().Name()
	}
	return res
}

func (p *property) Name() string {
	return p.name
}

func (p *property) TypeName() string {
	return p.typeName
}

func (p *property) Type() reflect.Type {
	return p.fieldInfo.Type
}

func (p *property) FieldInfo() reflect.StructField {
	return p.fieldInfo
}

func (p *property) Description() string {
	return p.description
}

func (p *property) Json() string {
	return p.json
}

func (p *property) IsArray() bool {
	return p.isArray
}

func (p *property) Kind() string {
	return p.kind
}
