package reflectutils

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
)

type Fields struct {
	items map[string]*reflect.StructField
}

func NewFields() *Fields {
	return &Fields{items: make(map[string]*reflect.StructField)}
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

type GetFieldsOptions struct {
	anonymous *bool // 是否包含匿名字段
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

func NewGetFieldsOptions(opts ...*GetFieldsOptions) *GetFieldsOptions {
	o := &GetFieldsOptions{anonymous: nil}
	for _, item := range opts {
		if item.anonymous != nil {
			o.anonymous = item.anonymous
		}
	}
	return o
}

// GetFields
//
//	@Description: 获取对象的字段信息
//	@param obj
//	@return *Fields
//	@return error
func GetFields(obj any, opts ...*GetFieldsOptions) (*Fields, error) {
	t := reflect.TypeOf(obj)
	return GetFieldsByType(t, opts...)
}

// GetFieldsByType
//
//	@Description:  根据对象类型,获取对象的字段信息
//	@param v
//	@return *Fields
//	@return error
func GetFieldsByType(v reflect.Type, opts ...*GetFieldsOptions) (*Fields, error) {
	o := NewGetFieldsOptions(opts...)
	fields := NewFields()
	t := v
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if _, err := addFieldsByType(t, fields, o.Anonymous()); err != nil {
		return nil, err
	}
	return fields, nil
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
