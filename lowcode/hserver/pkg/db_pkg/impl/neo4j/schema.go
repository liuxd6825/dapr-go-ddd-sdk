package neo4j

import (
	"context"
	dbschema "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/jsonschema/v6"
	"reflect"
	"time"
)

func newDbSchema(sch *jsonschema.Schema) *dbschema.Schema {
	s := dbschema.NewSchema()
	s.TableName = sch.Name
	s.Name = sch.Name
	props := sch.GetAllProperties()
	for _, prop := range props {
		dataType := getDataType(prop)
		f := &dbschema.Field{
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

func getSize(dataType dbschema.DataType, size *int) int {
	if size != nil {
		return *size
	}
	switch dataType {
	case dbschema.Int:
		return 8
	case dbschema.Float:
		return 8
	case dbschema.String:
		return 50
	case dbschema.Date:
		return 8
	case dbschema.Time:
		return 8
	case dbschema.Bool:
		return 1
	default:
		return 0
	}
}

func initField(field *dbschema.Field) {
	if field.DataType == dbschema.Time || field.DataType == dbschema.Date {
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
			} else if v, ok := val.(times.Time); ok {
				return v.PTime().UTC(), false
			} else if v, ok := val.(*times.Time); ok {
				return v.PTime().UTC(), false
			}
			panic("neither time nor timezone was set")
		}
	}
}

func getDataType(prop *jsonschema.Schema) dbschema.DataType {
	if prop.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return dbschema.Time
	} else if prop.Types.Contains(jsonschema.JsonType_StringType) {
		return dbschema.String
	} else if prop.Types.Contains(jsonschema.JsonType_IntegerType) {
		return dbschema.Int
	} else if prop.Types.Contains(jsonschema.JsonType_BooleanType) {
		return dbschema.Bool
	} else if prop.Types.Contains(jsonschema.JsonType_DateType) {
		return dbschema.Date
	}
	return dbschema.String
}
