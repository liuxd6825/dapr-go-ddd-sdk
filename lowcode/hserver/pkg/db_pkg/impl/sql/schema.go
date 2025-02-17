package sql

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	dbschema "gorm.io/gorm/schema"
	"reflect"
	"time"
)

// NewDBSchema get data type from dialector with extra schema table
func NewDBSchema(dest *schema.Schema) (*dbschema.Schema, error) {
	if dest == nil {
		return nil, fmt.Errorf("%w: %+v", dbschema.ErrUnsupportedDataType, dest)
	}

	s := newDbSchema(dest)
	var primaryField *dbschema.Field
	for _, p := range dest.Properties {
		var dataType dbschema.DataType
		name := stringutils.AsFieldName(p.Name)
		field := newDbField(s, name, dataType)
		field.NotNull = !p.IsTypeNull()
		setDbField(field, p.Field)
		if field.PrimaryKey {
			primaryField = field
		}
	}

	if primaryField == nil {
		idField := newDbField(s, "id", dbschema.String)
		idField.PrimaryKey = true
	}
	return s, nil
}

func newDbSchema(dest *schema.Schema) *dbschema.Schema {
	tableName := stringutils.AsFieldName(dest.Name)
	s := &dbschema.Schema{
		Name:             dest.Name,
		Table:            tableName,
		FieldsByName:     map[string]*dbschema.Field{},
		FieldsByBindName: map[string]*dbschema.Field{},
		FieldsByDBName:   map[string]*dbschema.Field{},
		Relationships:    dbschema.Relationships{Relations: map[string]*dbschema.Relationship{}},
	}
	return s
}

func newDbField(s *dbschema.Schema, name string, dataType dbschema.DataType) *dbschema.Field {
	fieldType := getFieldType(dataType)
	field := &dbschema.Field{
		Name:              name,
		DBName:            name,
		IndirectFieldType: fieldType,
		DataType:          dataType,
	}
	s.Fields = append(s.Fields, field)
	s.FieldsByDBName[name] = field
	s.FieldsByName[name] = field
	s.DBNames = append(s.DBNames, name)
	return field
}

func setDbField(field *dbschema.Field, propField *schema.Field) {
	if field == nil || propField == nil {
		return
	}

	field.PrimaryKey = propField.PrimaryKey
	field.NotNull = propField.NotNull
	field.DefaultValue = propField.DefaultValue
	if propField.Size != nil {
		field.Size = *propField.Size
	}
	field.Unique = propField.Unique
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
	default:
		return reflect.TypeOf("")
	}
}
