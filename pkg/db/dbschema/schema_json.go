package dbschema

import (
	"context"
	"reflect"
	"time"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	times2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	"github.com/liuxd6825/jsonschema/v6"
	gormschema "gorm.io/gorm/schema"
)

type DBSchema = store2.DBSchema

func NewDBSchemaWithJsonSchemaBytes(fileName string, jsonBytes []byte) *DBSchema {
	jsSchema := schema.NewJsonSchemaWithBytes(fileName, jsonBytes)
	return NewDBSchemaWithJsonSchema(jsSchema)
}

func NewDBSchemaWithJsonSchemaText(fileName string, jsonText string) *DBSchema {
	jsSchema := schema.NewJsonSchemaWithJson(fileName, jsonText)
	return NewDBSchemaWithJsonSchema(jsSchema)
}

func NewDBSchemaWithJsonSchema(sch *jsonschema.Schema) *DBSchema {
	s := store2.NewDBSchema()
	s.JsonSchema = sch
	s.TableName = schema.GetTableName(sch)
	s.Name = sch.Name()
	s.Title = sch.Title
	props := sch.GetAllProperties()
	for _, prop := range props {
		schField := schema.GetField(prop)
		if schField != nil && schField.NotField {
			continue
		}
		dataType := getDataType(prop)
		field := &store2.Field{
			Title:                 prop.Title,
			Name:                  prop.Name(),
			DBName:                AsFieldName(prop.Name()),
			DataType:              dataType,
			DefaultValueInterface: prop.Default,
			Creatable:             true,
			Updatable:             true,
			Readable:              true,
		}
		if schField != nil {
			if schField.Name != "" {
				field.DBName = schField.Name
			}
			field.Size = getSize(dataType, schField.Size)
			field.Updatable = schField.Updatable
			field.Creatable = schField.Creatable
			field.Readable = schField.Readable
			field.Unique = schField.Unique
			field.NotNull = schField.NotNull
			field.PrimaryKey = schField.PrimaryKey
			field.RelType = schField.RelType
			field.RelStartId = schField.RelStartId
			field.RelEndId = schField.RelEndId
			field.NodeLabel = schField.NodeLabel
			field.NodeLabelFormat = schField.NodeLabelFormat
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
	case gormschema.Int:
		return 8
	case gormschema.Float:
		return 8
	case gormschema.String:
		return 50
	case gormschema.Date:
		return 8
	case gormschema.Time:
		return 8
	case gormschema.Bool:
		return 1
	case gormschema.Json:
		return 100
	default:
		return 100
	}
}

func initField(field *store2.Field) {
	if field.DataType == gormschema.Time || field.DataType == gormschema.Date {
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
			} else if v, ok := val.(*times2.Date); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(times2.Date); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(times2.Time); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(*times2.Time); ok {
				return v.PTime().UTC(), false
			}
			panic("neither time nor timezone was set")
		}
	}
}

func getDataType(prop *jsonschema.Schema) DataType {
	if prop.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return gormschema.Time
	} else if prop.Types.Contains(jsonschema.JsonType_StringType) {
		return gormschema.String
	} else if prop.Types.Contains(jsonschema.JsonType_IntegerType) {
		return gormschema.Int
	} else if prop.Types.Contains(jsonschema.JsonType_BooleanType) {
		return gormschema.Bool
	} else if prop.Types.Contains(jsonschema.JsonType_DateType) {
		return gormschema.Date
	} else if prop.Types.Contains(jsonschema.JsonType_ObjectType) {
		return gormschema.Json
	} else if prop.Types.Contains(jsonschema.JsonType_ArrayType) {
		return gormschema.Array
	}
	return gormschema.String
}

func getTableName(sch *jsonschema.Schema) string {
	return schema.GetTableName(sch)
}
