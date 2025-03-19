package dbschema

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonschemautils"
	"github.com/liuxd6825/jsonschema/v6"
	"reflect"
	"time"
)

func NewDBSchemaWithJsonSchemaText(fileName string, jsonText string) *DBSchema {
	jsSchema := jsonschemautils.NewJsonSchemaWidthJson(fileName, jsonText)
	return NewDBSchemaWithJsonSchema(jsSchema)
}

func NewDBSchemaWithJsonSchema(sch *jsonschema.Schema) *DBSchema {
	s := NewDBSchema()
	s.TableName = sch.Name
	s.Name = sch.Name
	props := sch.GetAllProperties()
	for _, prop := range props {
		dataType := getDataType(prop)
		f := &Field{
			Name:                  prop.Name,
			DBName:                prop.DB.Name,
			DataType:              dataType,
			Size:                  getSize(dataType, prop.DB.Size),
			DefaultValueInterface: prop.Default,
			Updatable:             sch.DB.Updatable,
			Creatable:             sch.DB.Creatable,
			Readable:              sch.DB.Readable,
			Unique:                sch.DB.Unique,
			NotNull:               sch.DB.NotNull,
			PrimaryKey:            sch.DB.PrimaryKey,
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
