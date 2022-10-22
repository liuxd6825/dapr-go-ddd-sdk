package metadata

import (
	"fmt"
	"reflect"
)

type Properties interface {
	Values() map[string]Property
	Get(name string) (Property, bool)
}

type properties struct {
	values map[string]Property
}

type Options struct {
	Logger Logger
}

type Logger func(log string)

const (
	PropertyName       = "Property"
	PropertiesName     = "Properties"
	PkgPath            = "github.com/liuxd6825/dapr-go-ddd-sdk/metadata"
	DescriptionTagName = "description"
)

func NewOptions() *Options {
	o := &Options{}
	o.Logger = func(str string) {}
	return o
}

func NewProperties(metadata any, entity any, ops ...*Options) (Properties, error) {
	options := NewOptions().Merge(ops...)

	metaType := reflect.TypeOf(metadata)
	metaValue := reflect.ValueOf(metadata)

	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)

	if metaType.Kind() == reflect.Pointer {
		metaType = metaType.Elem()
		metaValue = metaValue.Elem()
	}

	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
		entityValue = entityValue.Elem()
	}

	if props, err := initProperties(metaType, metaValue, entityType, entityValue, options); err != nil {
		return nil, err
	} else {
		return props, nil
	}
}

func initProperties(metaType reflect.Type, metaValue reflect.Value, entityType reflect.Type, entityValue reflect.Value, options *Options) (*properties, error) {
	props := &properties{
		values: map[string]Property{},
	}
	// 取得对象的属性
	for i := 0; i < entityValue.NumField(); i++ {
		var prop *property
		entityField := entityType.Field(i)

		options.Logger(fmt.Sprintf("entityField.Name = %v", entityField.Name))
		if entityField.Anonymous {
			if mt, ok := metaType.FieldByName(entityField.Name + "Metadata"); ok {
				if et, ok := entityType.FieldByName(entityField.Name); ok {
					mv := metaValue.FieldByName(entityField.Name + "Metadata")
					ev := entityValue.FieldByName(entityField.Name)
					if sumProps, err := initProperties(mt.Type, mv, et.Type, ev, options); err != nil {
						return nil, err
					} else {
						props.AddProperties(sumProps)
					}
				}
			}
		} else if metaField, ok := metaType.FieldByName(entityField.Name); ok {
			t := metaField.Type
			if t.Name() == PropertyName && t.PkgPath() == PkgPath {
				fv := metaValue.FieldByName(metaField.Name)
				data := fv.Interface()
				if data == nil {
					prop = &property{}
					value := reflect.ValueOf(prop)
					fv.Set(value)
				} else if v, ok := data.(*property); ok {
					prop = v
				}
				if prop == nil {
					prop = &property{}
				}
				prop.Init(entityField)
				props.Add(prop)
			} else if metaField.Type.Name() == PropertiesName && metaField.Type.PkgPath() == PkgPath {
				v := reflect.ValueOf(props)
				fv := metaValue.FieldByName(metaField.Name)
				fv.Set(v)
			}
		}
	}
	return props, nil
}

func (m *properties) AddProperties(props Properties) {
	for _, item := range props.Values() {
		m.Add(item)
	}
}

func (m *properties) Values() map[string]Property {
	return m.values
}

func (m *properties) Add(p Property) {
	m.values[p.Name()] = p
}

func (m *properties) Get(name string) (Property, bool) {
	p, ok := m.values[name]
	return p, ok
}

func (o *Options) SetLogger(logger Logger) *Options {
	o.Logger = logger
	return o
}

func (o *Options) Merge(opts ...*Options) *Options {
	for _, i := range opts {
		if i.Logger != nil {
			o.Logger = i.Logger
		}
	}
	return o
}
