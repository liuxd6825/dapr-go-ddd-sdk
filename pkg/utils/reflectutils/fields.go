package reflectutils

import (
	"fmt"
	"reflect"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type Fields struct {
	items   map[string]*reflect.StructField
	refType reflect.Type
	refVal  reflect.Value
}

type GetFieldsOptions struct {
	anonymous *bool // 是否包含匿名字段
}

func newFields(obj any) *Fields {
	rType := reflect.TypeOf(obj)
	rVal := reflect.ValueOf(obj)

	fields := &Fields{
		items:   make(map[string]*reflect.StructField),
		refType: rType,
		refVal:  rVal,
	}
	return fields
}

// NewFields
//
//	@Description: 获取对象的字段信息
//	@param obj
//	@return *Fields
//	@return error
func NewFields(obj any, opts ...*GetFieldsOptions) (*Fields, error) {
	return NewFieldsByType(obj, opts...)
}

// NewFieldsByType
//
//	@Description:  根据对象类型,获取对象的字段信息
//	@param v
//	@return *Fields
//	@return error
func NewFieldsByType(obj any, opts ...*GetFieldsOptions) (*Fields, error) {
	o := NewGetFieldsOptions(opts...)
	fields := newFields(obj)
	t := fields.refType
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if _, err := addFieldsByType(t, fields, o.Anonymous()); err != nil {
		return nil, err
	}
	return fields, nil
}

func NewGetFieldsOptions(opts ...*GetFieldsOptions) *GetFieldsOptions {
	o := &GetFieldsOptions{anonymous: nil}
	for _, item := range opts {
		if item.anonymous != nil {
			o.anonymous = item.anonymous
		}
	}
	return o
}

func (f *Fields) addItem(v *reflect.StructField) {
	f.items[v.Name] = v
}

func (f *Fields) ForEach(foreach func(index int, name string, field *reflect.StructField)) {
	if foreach != nil {
		i := 0
		for name, item := range f.items {
			foreach(i, name, item)
			i++
		}
	}
}

func (f *Fields) Names() []string {
	var names []string
	for name, _ := range f.items {
		names = append(names, name)
	}
	return names
}

func (f *Fields) NamesFirstLower() []string {
	var names []string
	for name, _ := range f.items {
		names = append(names, stringutils.FirstLower(name))
	}
	return names
}

func (f *Fields) Contain(field string) bool {
	_, ok := f.items[field]
	return ok
}

func (f *Fields) Item(field string) (*reflect.StructField, bool) {
	i, ok := f.items[field]
	return i, ok
}

func (f *Fields) Count() int {
	return len(f.items)
}

func (f *Fields) SetValue(obj any, fieldName string, val any) {
	field, ok := f.Item(fieldName)
	if !ok {
		panic(fmt.Sprintf("field %s not found", fieldName))
	}

	refVal := reflect.ValueOf(obj)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(field.Name)

	// 根据字段类型设置值
	switch field.Type.Kind() {
	case reflect.String:
		if fieldValue.CanSet() {
			fieldValue.SetString(val.(string)) // 设置字符串字段的值
		}
	case reflect.Int:
		if fieldValue.CanSet() {
			fieldValue.SetInt(val.(int64)) // 设置整数字段的值
		}
	default:
		fmt.Println("Unsupported field type:", field.Type)
	}
}

func (f *Fields) GetString(obj any, fieldName string) string {
	field, ok := f.Item(fieldName)
	if !ok {
		panic(fmt.Sprintf("field %s not found", fieldName))
	}

	refVal := reflect.ValueOf(obj)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(field.Name)
	// 根据字段类型设置值
	switch field.Type.Kind() {
	case reflect.String:
		if fieldValue.CanSet() {
			return fieldValue.String() // 设置字符串字段的值
		}
	default:
		fmt.Println("Unsupported field type:", field.Type)
	}
	return ""
}

func (f *Fields) GetInt(obj any, fieldName string) int64 {
	field, ok := f.Item(fieldName)
	if !ok {
		panic(fmt.Sprintf("field %s not found", fieldName))
	}

	refVal := reflect.ValueOf(obj)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(field.Name)
	// 根据字段类型设置值
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() // 设置字符串字段的值
	default:
		fmt.Println("Unsupported field type:", field.Type)
	}
	return 0
}

func addFieldsByType(v reflect.Type, fields *Fields, anonymous bool) (*Fields, error) {
	t := v
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if anonymous {
				if _, err := addFieldsByType(f.Type, fields, anonymous); err != nil {
					return nil, err
				}
			}
		} else {
			fields.addItem(&f)
		}
	}

	return fields, nil
}

func (o *GetFieldsOptions) Anonymous() bool {
	if o.anonymous == nil {
		return true
	}
	return *o.anonymous
}

func (o *GetFieldsOptions) SetAnonymous(v bool) *GetFieldsOptions {
	o.anonymous = &v
	return o
}
