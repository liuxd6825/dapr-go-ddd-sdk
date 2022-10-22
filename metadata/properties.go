package metadata

import (
	"reflect"
	"testing"
)

type Properties interface {
	Values() map[string]Property
	Get(name string) (Property, bool)
}

type properties struct {
	values map[string]Property
}

const (
	PropertyName       = "Property"
	PkgPath            = "github.com/liuxd6825/dapr-go-ddd-sdk/metadata"
	DescriptionTagName = "description"
)

func NewProperties(metadata any, entity any, test *testing.T) (Properties, error) {
	metaType := reflect.TypeOf(metadata)
	metaValue := reflect.ValueOf(metadata)

	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)

	/*	if metaType.Kind() != reflect.Pointer {
			return nil, errors.New("metadata is not Pointer ")
		}
		if entityType.Kind() != reflect.Pointer {
			return nil, errors.New("entity is not Pointer ")
		}*/

	if metaType.Kind() == reflect.Pointer {
		metaType = metaType.Elem()
		metaValue = metaValue.Elem()
	}

	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
		entityValue = entityValue.Elem()
	}

	props := &properties{
		values: map[string]Property{},
	}

	if err := initMetadata(metaType, metaValue, entityType, entityValue, test, props); err != nil {
		return nil, err
	}
	return props, nil
}
func initMetadata(metaType reflect.Type, metaValue reflect.Value, entityType reflect.Type, entityValue reflect.Value, test *testing.T, props *properties) error {

	// 取得对象的属性
	for i := 0; i < entityValue.NumField(); i++ {
		var prop *property
		entityField := entityType.Field(i)
		test.Logf("entityField.Name = %v", entityField.Name)
		if entityField.Anonymous {
			if mt, ok := metaType.FieldByName(entityField.Name + "Metadata"); ok {
				mv := metaValue.FieldByName(entityField.Name + "Metadata")
				if et, ok := entityType.FieldByName(entityField.Name); ok {
					ev := entityValue.FieldByName(entityField.Name)
					if err := initMetadata(mt.Type, mv, et.Type, ev, test, props); err != nil {
						return err
					}
				}
			}
		}
		// 如果metadata对象中定义了相同字段
		if metaField, ok := metaType.FieldByName(entityField.Name); ok {
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
			}
		}
		if prop == nil {
			prop = &property{}
		}
		prop.Init(entityField)
		props.Add(prop)
	}
	return nil
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
