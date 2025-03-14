package dbschema

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/jsonschema/v6"
	"reflect"
	"time"
)

func NewSchemaWithJsonschema(sch *jsonschema.Schema) *Schema {
	s := NewSchema()
	s.TableName = sch.Name
	s.Name = sch.Name
	props := sch.GetAllProperties()
	for _, prop := range props {
		dataType := getDataType(prop)
		f := &Field{
			Name:                  prop.Name,
			DBName:                prop.Name,
			DataType:              dataType,
			Size:                  getSize(dataType, prop.DBField.Size),
			DefaultValueInterface: prop.Default,
			Updatable:             sch.DBField.Updatable,
			Creatable:             sch.DBField.Creatable,
			Readable:              sch.DBField.Readable,
			Unique:                sch.DBField.Unique,
			NotNull:               sch.DBField.NotNull,
			PrimaryKey:            sch.DBField.PrimaryKey,
		}
		initField(f)
		s.Fields = append(s.Fields, f)
	}
	return s
}

func getSize(dataType DataType, size *int) int {
	if size != nil {
		return *size
	}
	switch dataType {
	case Int:
		return 8
	case Float:
		return 8
	case String:
		return 50
	case Date:
		return 8
	case Time:
		return 8
	case Bool:
		return 1
	default:
		return 0
	}
}

func initField(field *Field) {
	if field.DataType == Time || field.DataType == Date {
		field.Set = func(ctx context.Context, value reflect.Value, i interface{}) error {
			value.Set(reflect.ValueOf(i))
			return nil
		}
		field.ValueOf = func(ctx context.Context, value reflect.Value) (res interface{}, zero bool) {
			val := value.Interface()
			if val == nil {
				return nil, true
			}
			if v, ok := val.(*time.Time); ok {
				return v.UTC(), false
			} else if v, ok := val.(time.Time); ok {
				return v.UTC(), false
			} else if v, ok := val.(*times.Date); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(times.Date); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(times.Time); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(*times.Time); ok {
				return v.PTime().UTC(), false
			}
			panic("neither time nor timezone was set")
		}
	}
}

func getDataType(prop *jsonschema.Schema) DataType {
	if prop.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return Time
	} else if prop.Types.Contains(jsonschema.JsonType_StringType) {
		return String
	} else if prop.Types.Contains(jsonschema.JsonType_IntegerType) {
		return Int
	} else if prop.Types.Contains(jsonschema.JsonType_BooleanType) {
		return Bool
	} else if prop.Types.Contains(jsonschema.JsonType_DateType) {
		return Date
	} else if prop.Types.Contains(jsonschema.JsonType_ObjectType) {
		return Object
	} else if prop.Types.Contains(jsonschema.JsonType_ArrayType) {
		return Array
	}
	return String
}
