package rsql

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
	"time"
)

func GetValue(value Value) interface{} {
	var v interface{}
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = sv.Value
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v = sv.Value
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v = sv.Value
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = GetValueList(sv)
	default:
		v = value
	}
	return v
}

func getValue(value Value) any {
	var v any
	var err error
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = sv.Value
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v, err = time.Parse(DateLayout, sv.Value)
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v, err = time.Parse(DateTimeLayout, sv.Value)
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = GetValueList(sv)
	case *FuncValue:
		sv, _ := value.(*FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

// AsFieldName
// @Description: 转换为mongodb规范的字段名称
// @param name
// @return string
func AsFieldName(name string) string {
	return stringutils.SnakeString(name)
}

func GetValueList(listValue *ListValue) []interface{} {
	list := make([]interface{}, 0)
	if listValue == nil {
		return list
	}
	for _, v := range listValue.Value {
		list = append(list, GetValue(v))
	}
	return list
}

// GetFieldValues 从 any 类型的列表中提取指定字段的值，并返回 []string
// 如果字段值是 string 类型，则添加单引号
func GetFieldValues(data any, field string) ([]string, error) {
	// 使用反射获取 data 的类型
	v := reflect.ValueOf(data)

	// 判断 data 是否是切片或数组
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("input is not a slice or array")
	}

	// 准备结果切片
	result := make([]string, 0, v.Len())

	// 遍历切片或数组中的每个元素
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		// 判断元素是结构体还是 map
		switch elem.Kind() {
		case reflect.Struct:
			// 如果是结构体，获取字段值
			fieldValue := elem.FieldByName(field)
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("field %s not found in struct", field)
			}
			result = append(result, getReflectValue(fieldValue))
		case reflect.Map:
			// 如果是 map，获取键对应的值
			fieldValue := elem.MapIndex(reflect.ValueOf(field))
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("field %s not found in map", field)
			}
			result = append(result, getReflectValue(fieldValue))
		default:
			return nil, fmt.Errorf("unsupported element type: %s", elem.Kind())
		}
	}
	return result, nil
}

// formatValue 根据字段值的类型格式化值
// 如果字段值是 string 类型，则添加单引号
func getReflectValue(value reflect.Value) string {
	// 获取字段值的实际类型
	switch value.Kind() {
	case reflect.String:
		// 如果是 string 类型，添加单引号
		return fmt.Sprintf("%v", value.Interface())
	default:
		// 其他类型直接转换为字符串
		return fmt.Sprintf("%v", value.Interface())
	}
}
