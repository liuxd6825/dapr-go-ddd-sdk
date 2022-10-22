package metadata

import "reflect"

type Property interface {
	Name() string
	TypeName() string
	Description() string
	JsonName() string
}

type property struct {
	FName        string              `json:"name"`
	FJsonName    string              `json:"jsonName"`
	FFieldInfo   reflect.StructField `json:"-"`
	FDescription string              `json:"description"`
}

func NewProperty() Property {
	return &property{}
}
func (p *property) Init(fieldInfo reflect.StructField) Property {
	p.FFieldInfo = fieldInfo
	p.FName = fieldInfo.Name
	if desc := fieldInfo.Tag.Get("description"); desc != "" {
		p.FDescription = desc
	}
	if jsonName := fieldInfo.Tag.Get("json"); jsonName != "" {
		p.FJsonName = jsonName
	}
	return p
}

func (p *property) Name() string {
	return p.FName
}

func (p *property) TypeName() string {
	return p.FFieldInfo.Type.Name()
}

func (p *property) Type() reflect.Type {
	return p.FFieldInfo.Type
}

func (p *property) FieldInfo() reflect.StructField {
	return p.FFieldInfo
}

func (p *property) Description() string {
	return p.FDescription
}

func (p *property) JsonName() string {
	return p.FJsonName
}
