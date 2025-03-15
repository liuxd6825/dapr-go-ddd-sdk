package dbschema

import (
	"context"
	"fmt"
	"reflect"
)

type Schema struct {
	Name        string
	TableName   string
	Fields      []*Field
	fieldName   map[string]*Field
	fieldDbName map[string]*Field
}

func NewSchema() *Schema {
	return &Schema{
		Name:        "",
		TableName:   "",
		Fields:      []*Field{},
		fieldName:   map[string]*Field{},
		fieldDbName: map[string]*Field{},
	}
}

func (sch *Schema) SetName(name string) *Schema {
	sch.Name = name
	return sch
}

func (sch *Schema) SetTableName(name string) *Schema {
	sch.TableName = name
	return sch
}

func (sch *Schema) AddField(field ...*Field) *Schema {
	sch.Fields = append(sch.Fields, field...)
	for _, f := range sch.Fields {
		sch.fieldName[f.Name] = f
		sch.fieldDbName[f.DBName] = f
	}
	return sch
}

func (sch *Schema) InitFields() {
	sch.fieldName = nil
	sch.initFields()
}

func (sch *Schema) initFields() {
	if sch.fieldName == nil {
		sch.fieldName = map[string]*Field{}
		sch.fieldDbName = map[string]*Field{}
	} else {
		return
	}
	for _, field := range sch.Fields {
		sch.fieldName[field.Name] = field
		sch.fieldDbName[field.DBName] = field
	}
}

func (sch *Schema) LookedField(name string) *Field {
	if f, ok := sch.fieldDbName[name]; ok {
		return f
	}
	if f, ok := sch.fieldName[name]; ok {
		return f
	}
	return nil
}

func (sch *Schema) NewMap(ctx context.Context, obj any, opts ...func(map[string]any)) (map[string]any, error) {
	sch.initFields()

	res := map[string]any{}
	for _, opt := range opts {
		opt(res)
	}

	if obj == nil {
		return res, nil
	}
	if e, ok := any(obj).(map[string]any); ok {
		for _, field := range sch.Fields {
			if field == nil {
				continue
			}
			val, ok := e[field.Name]
			if field.ValueOf != nil && ok {
				val, _ = field.ValueOf(ctx, reflect.ValueOf(val))
			}
			res[field.Name] = val
		}
	} else {
		vObj := reflect.ValueOf(obj)
		if vObj.Kind() == reflect.Ptr {
			vObj = vObj.Elem()
		}

		if vObj.Kind() != reflect.Struct {
			return nil, fmt.Errorf("input is not a struct")
		}

		for _, field := range sch.Fields {
			fv := vObj.FieldByName(field.Name)
			val := fv.Interface()
			if fv.Kind() == reflect.Struct {
				nestedMap, err := sch.NewMap(ctx, val)
				if err != nil {
					return nil, err
				}
				res[field.Name] = nestedMap
			} else if field.ValueOf != nil {
				v, _ := field.ValueOf(ctx, fv)
				val = v
			}
			if field.ValueOf != nil {
				val, _ = field.ValueOf(ctx, reflect.ValueOf(val))
			}
			res[field.Name] = val
		}
	}
	return res, nil
}
