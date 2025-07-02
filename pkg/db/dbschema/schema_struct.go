package dbschema

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	gormschema "gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewDBSchema(name string, tableName string) *store.DBSchema {
	var err error = nil
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			panic(fmt.Errorf("tableName:%s; %s", tableName, err.Error()))
		}
	}()
	dbSch := store.NewDBSchema()
	dbSch.SetTableName(tableName)
	dbSch.SetName(name)
	return dbSch
}

func NewDBSchemaWithStruct(name string, data any, tableName string) *store.DBSchema {
	var err error = nil
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			panic(fmt.Errorf("tableName:%s; %s", tableName, err.Error()))
		}
	}()
	gSch, err := gormschema.ParseWithSpecialTableName(data, &sync.Map{}, gormschema.NamingStrategy{}, tableName)
	if err != nil {
		panic(err)
	}
	dbSch := store.NewDBSchema()
	for _, f := range gSch.Fields {
		field := store.NewField()
		field.Name = f.Name
		field.DBName = f.DBName
		field.StructField = f.StructField
		field.IndirectFieldType = f.IndirectFieldType
		field.Serializer = f.Serializer
		field.FieldType = f.FieldType
		field.Tag = f.Tag
		field.FieldType = f.FieldType
		field.DataType = f.DataType
		field.ValueOf = f.ValueOf
		field.Creatable = f.Creatable
		field.Updatable = f.Updatable
		field.PrimaryKey = f.PrimaryKey

		field.RelType = f.RelType
		field.RelEndId = f.RelEndId
		field.RelStartId = f.RelStartId
		field.NodeLabelFormat = f.NodeLabelFormat
		field.NodeLabel = f.NodeLabel

		dbSch.AddField(field)
	}
	dbSch.SetTableName(tableName)
	dbSch.SetName(name)
	dbSch.GormSchema = gSch
	return dbSch
}

func NewGormSchema(dbSch *store.DBSchema) *gormschema.Schema {
	var err error = nil
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			panic(fmt.Errorf("tableName:%s; %s", dbSch.TableName, err.Error()))
		}
	}()
	if dbSch == nil {
		panic("NewGormSchema() dbSch=nil")
	}
	gSch := gormschema.NewSchemaEmpty()
	for _, f := range dbSch.Fields {
		field := gormschema.NewField()
		field.Name = f.Name
		field.DBName = f.DBName
		field.StructField = f.StructField
		field.IndirectFieldType = f.IndirectFieldType
		field.Serializer = f.Serializer
		field.FieldType = f.FieldType
		field.Tag = f.Tag
		field.FieldType = f.FieldType
		field.DataType = f.DataType
		field.ValueOf = f.ValueOf
		field.Creatable = f.Creatable
		field.Updatable = f.Updatable

		field.RelType = f.RelType
		field.RelEndId = f.RelEndId
		field.RelStartId = f.RelStartId
		field.NodeLabelFormat = f.NodeLabelFormat
		field.NodeLabel = f.NodeLabel

		gSch.AddField(field)
	}
	gSch.Table = dbSch.TableName
	gSch.Name = dbSch.Name
	return gSch
}

func getFieldsByValue(refVal reflect.Value) []*store.Field {
	// 如果传入的是指针，解引用
	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	fields := make([]*store.Field, 0)
	for i := 0; i < refVal.NumField(); i++ {
		rField := refVal.Type().Field(i)
		rVal := refVal.Field(i)

		println("field.Name=", rField.Name)
		if rField.Name == "Birthday" {
			println(rField.Tag.Get("Birthday"))
		}
		tagSetting := ParseTagSetting(rField.Tag.Get("gorm"), ";")
		// 如果字段是嵌入结构体，递归处理
		if rField.Anonymous {
			fields = append(fields, getFieldsByValue(rVal)...)
		}
		field := &store.Field{
			DBName:            rField.Name,
			Name:              rField.Name,
			FieldType:         rField.Type,
			IndirectFieldType: rField.Type,
			StructField:       rField,
			Tag:               rField.Tag,
			TagSettings:       tagSetting,
		}
		setFieldByReflect(rVal, field)
		fields = append(fields, field)

	}
	return fields
}

func setFieldByReflect(fieldValue reflect.Value, field *store.Field) {
	var err error
	skipParseDefaultValue := strings.Contains(field.DefaultValue, "(") &&
		strings.Contains(field.DefaultValue, ")") || strings.ToLower(field.DefaultValue) == "null" || field.DefaultValue == ""

	switch reflect.Indirect(fieldValue).Kind() {
	case reflect.Bool:
		field.DataType = gormschema.Bool
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseBool(field.DefaultValue); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for bool, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.DataType = gormschema.Int
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseInt(field.DefaultValue, 0, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for int, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.DataType = gormschema.Uint
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseUint(field.DefaultValue, 0, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for uint, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Float32, reflect.Float64:
		field.DataType = gormschema.Float
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseFloat(field.DefaultValue, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for float, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.String:
		field.DataType = gormschema.String
		if field.HasDefaultValue && !skipParseDefaultValue {
			field.DefaultValue = strings.Trim(field.DefaultValue, "'")
			field.DefaultValue = strings.Trim(field.DefaultValue, `"`)
			field.DefaultValueInterface = field.DefaultValue
		}
	case reflect.Struct:
		if _, ok := fieldValue.Interface().(*time.Time); ok {
			field.DataType = gormschema.Time
		} else if fieldValue.Type().ConvertibleTo(store.TimeReflectType) {
			field.DataType = gormschema.Time
		} else if fieldValue.Type().ConvertibleTo(store.TimePtrReflectType) {
			field.DataType = gormschema.Time
		}
		if field.HasDefaultValue && !skipParseDefaultValue && (field.DataType == gormschema.Time || field.DataType == gormschema.Date) {
			if t, err := now.Parse(field.DefaultValue); err == nil {
				field.DefaultValueInterface = t
			}
		}
	case reflect.Array, reflect.Slice:
		if reflect.Indirect(fieldValue).Type().Elem() == store.ByteReflectType && field.DataType == "" {
			field.DataType = gormschema.Bytes
		} else {
			field.DataType = gormschema.Json
		}
	}

	if dataTyper, ok := fieldValue.Interface().(gormschema.GormDataTypeInterface); ok {
		field.DataType = DataType(dataTyper.GormDataType())
	}
	field.DataType = "string"

	if field.DataType == "" {
		panic(fmt.Errorf("failed to parse %s as default value for struct field, got nil", field.Name))
	}
}
