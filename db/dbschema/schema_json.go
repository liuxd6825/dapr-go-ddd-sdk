package dbschema

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonschemautils"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/liuxd6825/jsonschema/v6/extensions"
	"reflect"
	"time"
)

func NewDBSchemaWithJsonSchemaText(fileName string, jsonText string) *DBSchema {
	jsSchema := jsonschemautils.NewJsonSchemaWidthJson(fileName, jsonText)
	return NewDBSchemaWithJsonSchema(jsSchema)
}

func NewDBSchemaWithJsonSchema(sch *jsonschema.Schema) *DBSchema {
	s := NewDBSchema()
	s.TableName = extensions.GetTableName(sch)
	s.Name = sch.Name()
	props := sch.GetAllProperties()

	for _, prop := range props {
		metaField := extensions.GetField(prop)
		dataType := getDataType(prop)
		field := &Field{
			Name:                  prop.Name(),
			DataType:              dataType,
			DefaultValueInterface: prop.Default,
		}
		if metaField != nil {
			field.DBName = extensions.GetFieldName(prop)
			field.Size = getSize(dataType, metaField.Size)
			field.Updatable = metaField.Updatable
			field.Creatable = metaField.Creatable
			field.Readable = metaField.Readable
			field.Unique = metaField.Unique
			field.NotNull = metaField.NotNull
			field.PrimaryKey = metaField.PrimaryKey
		}
		initField(field)
		s.Fields = append(s.Fields, field)
	}
	return s
}

func getSize(dataType DataType, size int64) int {
	if size != 0 {
		return int(size)
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

func getTableName(sch *jsonschema.Schema) string {
	return extensions.GetTableName(sch)
}
