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
func NewDBSchema(sch *jsonschema.Schema) (*dbschema.Schema, error) {
	if sch == nil {
		return nil, fmt.Errorf("%w: %+v", dbschema.ErrUnsupportedDataType, sch)
	}

	dbSch := newDbSchema(sch)
	primaryField := addFields(sch, dbSch)

	if primaryField == nil {
		idField := dbSch.FieldsByDBName["id"]
		if idField == nil {
			idField = addDbField(dbSch, "id", dbschema.String, 30)
		}
		idField.PrimaryKey = true
	}

	tenantIdField := dbSch.FieldsByName["tenantId"]
	if tenantIdField == nil {
		tenantIdField = addDbField(dbSch, "tenantId", dbschema.String, 20)
	} else {
		tenantIdField.NotNull = true
	}

	initDbSchema(dbSch)
	return dbSch, nil
}

func addFields(sch *jsonschema.Schema, dbSch *dbschema.Schema) (primaryField *dbschema.Field) {
	for _, p := range sch.AllOf {
		if p.Ref != nil {
			if v := addFields(p.Ref, dbSch); v != nil {
				primaryField = v
			}
		}
	}
	for _, p := range sch.Properties {
		dataType := getDataType(p)
		size := 0
		field := addDbField(dbSch, p.Name, dataType, size)
		if field.PrimaryKey {
			primaryField = field
		}
		setDbField(field, p)
	}
	return primaryField
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

	//addDbField(s, "tenantId", dbschema.String)
	return s
}
func getDataType(property *jsonschema.Schema) dbschema.DataType {
	if property == nil {
		panic("getDataType: property is nil")
	}

	if property.Types.Contains(jsonschema.JsonType_DateType) {
		return dbschema.Time
	}
	if property.Types.Contains(jsonschema.JsonType_DateTimeType) {
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
func addDbField(s *dbschema.Schema, name string, dataType dbschema.DataType, size int) *dbschema.Field {
	dbName := stringutils.AsFieldName(name)
	fieldType := getFieldType(dataType)
	fieldSize := getFieldSize(dataType, size)
	field := &dbschema.Field{
		FieldType:         fieldType,
		IndirectFieldType: fieldType,
		Name:              name,
		DBName:            dbName,
		DataType:          dataType,
		Creatable:         true,
		Updatable:         true,
		Readable:          true,
		Size:              fieldSize,
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

func getFieldSize(fieldType dbschema.DataType, size int) int {
	if size > 0 {
		return size
	}
	switch fieldType {
	case dbschema.Bool:
		return 2
	case dbschema.String:
		return 100
	case dbschema.Int:
		return 10
	case dbschema.Float:
		return 10
	case dbschema.Time:
		return 10
	case dbschema.Object:
		return 1000
	case dbschema.Array:
		return 100
	default:
		return size
	}
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
