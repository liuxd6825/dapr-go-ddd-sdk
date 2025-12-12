package store

import (
	"context"
	"fmt"
	"reflect"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	gormschema "gorm.io/gorm/schema"
)

type DBSchema struct {
	Name              string
	TableName         string
	Title             string
	Fields            []*Field
	FieldName         map[string]*Field
	FieldDbName       map[string]*Field
	Labels            []string
	GormSchema        *gormschema.Schema
	JsonSchema        *jsonschema.Schema
	relTypeField      *Field
	relTypeFieldOk    bool
	relStartIdField   *Field
	relStartIdFieldOk bool
	relEndIdField     *Field
	relEndIdFieldOk   bool
	labelFields       []*Field
	labelFieldsOK     bool
	/*	nameFieldsOK      bool
		nameFields        []*Field*/
}

func NewDBSchema() *DBSchema {
	return &DBSchema{
		Name:        "",
		TableName:   "",
		Fields:      []*Field{},
		FieldName:   map[string]*Field{},
		FieldDbName: map[string]*Field{},
	}
}

func (sch *DBSchema) SetName(name string) *DBSchema {
	sch.Name = name
	return sch
}

func (sch *DBSchema) SetTableName(name string) *DBSchema {
	sch.TableName = name
	return sch
}

func (sch *DBSchema) AddField(field ...*Field) *DBSchema {
	sch.Fields = append(sch.Fields, field...)
	for _, f := range sch.Fields {
		sch.FieldName[f.Name] = f
		sch.FieldDbName[f.DBName] = f
	}
	return sch
}

func (sch *DBSchema) InitFields() {
	sch.initFields()
}

func (sch *DBSchema) initFields() {
	if len(sch.FieldName) == 0 {
		for _, field := range sch.Fields {
			sch.FieldName[field.Name] = field
			sch.FieldDbName[field.DBName] = field
		}
	}
}

func (sch *DBSchema) LookedField(name string) *Field {
	sch.initFields()
	name = stringutils.AsFieldName(name)
	if f, ok := sch.FieldDbName[name]; ok {
		return f
	}
	if f, ok := sch.FieldName[name]; ok {
		return f
	}
	return nil
}

func (sch *DBSchema) GetNodeLabelFields() []*Field {
	if !sch.labelFieldsOK {
		var labelFields []*Field
		for _, field := range sch.Fields {
			if field.NodeLabel {
				labelFields = append(labelFields, field)
			}
		}
		sch.labelFields = labelFields
		sch.labelFieldsOK = true
	}
	return sch.labelFields
}

func (sch *DBSchema) GetRelTypeField() *Field {
	if !sch.relTypeFieldOk {
		for _, field := range sch.Fields {
			if field.RelType {
				sch.relTypeField = field
				break
			}
		}
		sch.relTypeFieldOk = true
	}
	return sch.relTypeField
}

func (sch *DBSchema) GetRelStartIdField() *Field {
	if !sch.relStartIdFieldOk {
		for _, field := range sch.Fields {
			if field.RelStartId {
				sch.relStartIdField = field
				break
			}
		}
		sch.relStartIdFieldOk = true
	}
	return sch.relStartIdField
}

func (sch *DBSchema) GetRelEndIdField() *Field {
	if !sch.relEndIdFieldOk {
		for _, field := range sch.Fields {
			if field.RelEndId {
				sch.relEndIdField = field
				break
			}
		}
		sch.relEndIdFieldOk = true
	}
	return sch.relEndIdField
}

func (sch *DBSchema) NewMap(ctx context.Context, obj any, opts ...func(map[string]any)) (res map[string]any, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			println(err.Error())
		}
	}()

	sch.initFields()

	res = map[string]any{}
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
			key := field.Name
			val, ok := e[key]

			if ok {
				if field.ValueOf != nil && val != nil {
					//val, _ = field.ValueOf(ctx, reflect.ValueOf(val))
				}
			} else {
				val = nil
			}
			res[key] = val
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
			key := field.Name
			fv := vObj.FieldByName(key)
			var val any
			if fv.IsValid() {
				val = fv.Interface()
			}
			if fv.Kind() == reflect.Struct {
				nestedMap, err := sch.NewMap(ctx, val)
				if err != nil {
					return nil, err
				}
				res[key] = nestedMap
			}
			/*else if field.ValueOf != nil {
				v, _ := field.ValueOf(ctx, fv)
				val = v
			}
			if field.ValueOf != nil {
				val, _ = field.ValueOf(ctx, reflect.ValueOf(val))
			}*/
			res[key] = val
		}
	}
	return res, nil
}
