package metadata

import (
	"reflect"
)

type Properties interface {
	Values() map[string]Property
	Add(p Property) error
	Get(name string) (Property, bool)
}

type properties struct {
	values map[string]Property
}

const (
	PropertyName = "Property"
	PkgPath      = "github.com/liuxd6825/dapr-go-ddd-sdk/metadata"
)

func NewProperties(obj interface{}) Properties {
	m := &properties{
		values: map[string]Property{},
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if ft.Type.Name() == PropertyName && ft.Type.PkgPath() == PkgPath {
			fv := v.FieldByName(ft.Name)
			p := fv.Interface().(Property)
			m.Add(p)
		}
	}
	return m
}

func (m *properties) Values() map[string]Property {
	return m.values
}

func (m *properties) Add(p Property) error {
	m.values[p.Name()] = p
	return nil
}

func (m *properties) Get(name string) (Property, bool) {
	p, ok := m.values[name]
	return p, ok
}
