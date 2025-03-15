package dbschema

import (
	"fmt"
	"github.com/jinzhu/now"
	gormschema "gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewSchemaWithStruct(name string, data any, tableName string) *Schema {
	sch := NewSchema()
	sch.SetTableName(tableName)
	sch.SetName(name)
	v := reflect.ValueOf(data)
	fields := getFieldsByValue(v)
	for _, field := range fields {
		sch.AddField(field)
	}
	return sch
}

func getFieldsByValue(refVal reflect.Value) []*Field {
	// 如果传入的是指针，解引用
	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	fields := make([]*Field, 0)
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
		field := &Field{
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

func setFieldByReflect(fieldValue reflect.Value, field *Field) {
	var err error
	skipParseDefaultValue := strings.Contains(field.DefaultValue, "(") &&
		strings.Contains(field.DefaultValue, ")") || strings.ToLower(field.DefaultValue) == "null" || field.DefaultValue == ""

	switch reflect.Indirect(fieldValue).Kind() {
	case reflect.Bool:
		field.DataType = Bool
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseBool(field.DefaultValue); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for bool, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.DataType = Int
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseInt(field.DefaultValue, 0, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for int, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.DataType = Uint
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseUint(field.DefaultValue, 0, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for uint, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.Float32, reflect.Float64:
		field.DataType = Float
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseFloat(field.DefaultValue, 64); err != nil {
				panic(fmt.Errorf("failed to parse %s as default value for float, got error: %v", field.DefaultValue, err))
			}
		}
	case reflect.String:
		field.DataType = String
		if field.HasDefaultValue && !skipParseDefaultValue {
			field.DefaultValue = strings.Trim(field.DefaultValue, "'")
			field.DefaultValue = strings.Trim(field.DefaultValue, `"`)
			field.DefaultValueInterface = field.DefaultValue
		}
	case reflect.Struct:
		if _, ok := fieldValue.Interface().(*time.Time); ok {
			field.DataType = Time
		} else if fieldValue.Type().ConvertibleTo(TimeReflectType) {
			field.DataType = Time
		} else if fieldValue.Type().ConvertibleTo(TimePtrReflectType) {
			field.DataType = Time
		}
		if field.HasDefaultValue && !skipParseDefaultValue && (field.DataType == Time || field.DataType == Date) {
			if t, err := now.Parse(field.DefaultValue); err == nil {
				field.DefaultValueInterface = t
			}
		}
	case reflect.Array, reflect.Slice:
		if reflect.Indirect(fieldValue).Type().Elem() == ByteReflectType && field.DataType == "" {
			field.DataType = Bytes
		} else {
			field.DataType = Array
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
