package sql

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	dbschema "gorm.io/gorm/schema"
	"reflect"
	"time"
)

// NewDBSchema get data type from dialector with extra schema table
func NewDBSchema(dest *jsonschema.Schema) (*dbschema.Schema, error) {
	if dest == nil {
		return nil, fmt.Errorf("%w: %+v", dbschema.ErrUnsupportedDataType, dest)
	}

	s := newDbSchema(dest)
	var primaryField *dbschema.Field
	for _, p := range dest.Properties {
		dataType := getDataType(p)
		field := addDbField(s, p.Name, dataType)
		setDbField(field, p)
		if field.PrimaryKey {
			primaryField = field
		}
	}

	if primaryField == nil {
		idField := addDbField(s, "id", dbschema.String)
		idField.PrimaryKey = true
	}
	initDbSchema(s)
	return s, nil
}

func initDbSchema(s *dbschema.Schema) {

	for k, f := range s.FieldsByDBName {
		if f.PrimaryKey {
			s.PrimaryFields = append(s.PrimaryFields, f)
			s.PrimaryFieldDBNames = append(s.PrimaryFieldDBNames, k)
		}
	}
}

func newDbSchema(dest *jsonschema.Schema) *dbschema.Schema {
	tableName := stringutils.AsFieldName(dest.Name)
	s := &dbschema.Schema{
		Name:                dest.Name,
		Table:               tableName,
		DBNames:             []string{},
		FieldsByName:        map[string]*dbschema.Field{},
		FieldsByBindName:    map[string]*dbschema.Field{},
		FieldsByDBName:      map[string]*dbschema.Field{},
		Relationships:       dbschema.Relationships{Relations: map[string]*dbschema.Relationship{}},
		PrimaryFields:       make([]*dbschema.Field, 0),
		PrimaryFieldDBNames: []string{},
		ModelType:           reflect.TypeOf(map[string]any{}),
	}

	addDbField(s, "tenantId", dbschema.String)
	return s
}
func getDataType(property *jsonschema.Schema) dbschema.DataType {
	if property == nil {
		panic("getDataType: property is nil")
	}

	if property.Types.Contains(jsonschema.JsonType_DateType) {
		return dbschema.Time
	}
	if property.Types.Contains(jsonschema.JsonType_IntegerType) {
		return dbschema.Int
	}
	if property.Types.Contains(jsonschema.JsonType_BooleanType) {
		return dbschema.Bool
	}
	if property.Types.Contains(jsonschema.JsonType_StringType) {
		return dbschema.String
	}
	if property.Types.Contains(jsonschema.JsonType_ObjectType) {
		return dbschema.Object
	}
	if property.Types.Contains(jsonschema.JsonType_ArrayType) {
		return dbschema.Array
	}
	return dbschema.String
}
func addDbField(s *dbschema.Schema, name string, dataType dbschema.DataType) *dbschema.Field {
	dbName := stringutils.AsFieldName(name)
	fieldType := getFieldType(dataType)
	field := &dbschema.Field{
		FieldType:         fieldType,
		IndirectFieldType: fieldType,
		Name:              name,
		DBName:            dbName,
		DataType:          dataType,
		Creatable:         true,
		Updatable:         true,
		Readable:          true,
	}
	s.Fields = append(s.Fields, field)
	s.FieldsByDBName[dbName] = field
	s.FieldsByName[name] = field
	s.DBNames = append(s.DBNames, dbName)
	return field
}

func setDbField(field *dbschema.Field, propField *jsonschema.Schema) {
	if field == nil || propField == nil {
		return
	}
	/*
		field.PrimaryKey = propField.PrimaryKey
		field.NotNull = propField.NotNull
		field.DefaultValue = propField.DefaultValue
		if propField.Size != nil {
			field.Size = *propField.Size
		}
		field.Unique = propField.Unique

	*/
}

func getFieldType(dbType dbschema.DataType) reflect.Type {
	switch dbType {
	case dbschema.Bool:
		return reflect.TypeOf(true)
	case dbschema.String:
		return reflect.TypeOf("")
	case dbschema.Int:
		return reflect.TypeOf(int(0))
	case dbschema.Float:
		return reflect.TypeOf(float64(0))
	case dbschema.Time:
		return reflect.TypeOf(time.Time{})
	case dbschema.Object:
		return reflect.TypeOf("map[string]any{}")
	case dbschema.Array:
		return reflect.TypeOf("[]any{}")
	default:
		return reflect.TypeOf("")
	}
}
