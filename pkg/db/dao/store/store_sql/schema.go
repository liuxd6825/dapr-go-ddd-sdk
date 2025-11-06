package store_sql

import (
	"fmt"
	"reflect"
	"time"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
	gormschema "gorm.io/gorm/schema"
)

// NewGormSchema get data type from dialector with extra schema table
func NewGormSchema(sch *store2.DBSchema) (*gormschema.Schema, error) {
	if sch == nil {
		return nil, fmt.Errorf("%w: %+v", gormschema.ErrUnsupportedDataType, sch)
	}

	gormSch := newGormSchema(sch)
	primaryField := addFields(sch, gormSch)

	if primaryField == nil {
		idField := gormSch.FieldsByDBName["id"]
		if idField == nil {
			idField = addDbField(gormSch, "id", gormschema.String, 30)
		}
		idField.PrimaryKey = true
		idField.Updatable = false
	}

	tenantIdField := gormSch.FieldsByName["tenantId"]
	if tenantIdField == nil {
		tenantIdField = addDbField(gormSch, "tenantId", gormschema.String, 20)
	} else {
		tenantIdField.NotNull = true
	}
	tenantIdField.Updatable = false

	initGormSchema(gormSch)
	return gormSch, nil
}

func addFields(sch *store2.DBSchema, dbSch *gormschema.Schema) (primaryField *gormschema.Field) {
	for _, f := range sch.Fields {
		dataType := getDataType(f)
		size := 0
		field := addDbField(dbSch, f.Name, dataType, size)
		if field.PrimaryKey {
			primaryField = field
		}
		setDbField(field, f)
	}
	return primaryField
}

func initGormSchema(s *gormschema.Schema) {
	for k, f := range s.FieldsByDBName {
		if f.PrimaryKey {
			s.PrimaryFields = append(s.PrimaryFields, f)
			s.PrimaryFieldDBNames = append(s.PrimaryFieldDBNames, k)
		}
	}
}

func newGormSchema(dest *store2.DBSchema) *gormschema.Schema {
	tableName := stringutils.AsFieldName(dest.Name)
	s := &gormschema.Schema{
		Name:                dest.Name,
		Table:               tableName,
		DBNames:             []string{},
		FieldsByName:        map[string]*gormschema.Field{},
		FieldsByBindName:    map[string]*gormschema.Field{},
		FieldsByDBName:      map[string]*gormschema.Field{},
		Relationships:       gormschema.Relationships{Relations: map[string]*gormschema.Relationship{}},
		PrimaryFields:       make([]*gormschema.Field, 0),
		PrimaryFieldDBNames: []string{},
		ModelType:           reflect.TypeOf(map[string]any{}),
	}

	//addDbField(s, "tenantId", dbschema.String)
	return s
}
func getDataType(field *store2.Field) gormschema.DataType {
	if field == nil {
		panic("getDataType: property is nil")
	}

	switch field.DataType {
	case store2.DataType_Date:
		return gormschema.Date
	case store2.DataType_Time:
		return gormschema.Time
	case store2.DataType_String:
		return gormschema.String
	case store2.DataType_Int:
		return gormschema.Int
	case store2.DataType_Float:
		return gormschema.Float
	case store2.DataType_Bool:
		return gormschema.Bool
	case store2.DataType_Bytes:
		return gormschema.Bytes
	case store2.DataType_Array:
		return gormschema.Json
	case store2.DataType_Json:
		return gormschema.Json
	case store2.DataType_Uint:
		return gormschema.Uint
	}
	return gormschema.String
}
func addDbField(s *gormschema.Schema, name string, dataType gormschema.DataType, size int) *gormschema.Field {
	dbName := stringutils.AsFieldName(name)
	fieldType := getFieldType(dataType)
	fieldSize := getFieldSize(dataType, size)
	field := &gormschema.Field{
		FieldType:         fieldType,
		IndirectFieldType: fieldType,
		Name:              name,
		DBName:            dbName,
		DataType:          dataType,
		Creatable:         true,
		Updatable:         true,
		Readable:          true,
		Size:              fieldSize,
		TagSettings:       nil,
	}
	s.Fields = append(s.Fields, field)
	s.FieldsByDBName[dbName] = field
	s.FieldsByName[name] = field
	s.DBNames = append(s.DBNames, dbName)
	return field
}

func setDbField(field *gormschema.Field, propField *store2.Field) {
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

func getFieldSize(fieldType gormschema.DataType, size int) int {
	if size > 0 {
		return size
	}
	switch fieldType {
	case gormschema.Bool:
		return 2
	case gormschema.String:
		return 100
	case gormschema.Int:
		return 10
	case gormschema.Float:
		return 10
	case gormschema.Time:
		return 10
	case gormschema.Date:
		return 10
	case gormschema.Json:
		return 1000
	default:
		return size
	}
}

func getFieldType(dbType gormschema.DataType) reflect.Type {
	switch dbType {
	case gormschema.Bool:
		return reflect.TypeOf(true)
	case gormschema.String:
		return reflect.TypeOf("")
	case gormschema.Int:
		return reflect.TypeOf(int(0))
	case gormschema.Float:
		return reflect.TypeOf(float64(0))
	case gormschema.Time:
		return reflect.TypeOf(time.Time{})
	case gormschema.Date:
		return reflect.TypeOf(time.Time{})
	case gormschema.Json:
		return reflect.TypeOf("map[string]any{}")
	default:
		return reflect.TypeOf("")
	}
}
