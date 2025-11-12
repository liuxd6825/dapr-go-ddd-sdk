package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

type ConvertType string
type String interface {
	String() string
}

const (
	ConvertTypeNumber   ConvertType = "number"
	ConvertTypeInt      ConvertType = "int"
	ConvertTypeFloat    ConvertType = "float"
	ConvertTypeString   ConvertType = "string"
	ConvertTypeBool     ConvertType = "bool"
	ConvertTypeDateTime ConvertType = "datetime"
	ConvertTypeTime     ConvertType = "time"
	ConvertTypeNone     ConvertType = ""
)

func Convert(convType ConvertType, value any) (any, error) {
	switch convType {
	case ConvertTypeNumber, ConvertTypeInt:
		return ToInt(value)
	case ConvertTypeFloat:
		return ToFloat64(value)
	case ConvertTypeString,
		ConvertTypeNone:
		return ToString(value)
	case ConvertTypeBool:
		return ToBool(value)
	case ConvertTypeDateTime:
		return ToDateTime(value)
	case ConvertTypeTime:
		return ToTime(value)
	default:
		typeName := reflect.ValueOf(value).Type().String()
		return nil, errors.New(fmt.Sprintf("unknown convert type: %s , %s", convType, typeName))
	}
}

func ToString(value any) (string, error) {
	v := convertor.ToString(value)
	return v, nil
}

func ToInt(value any) (res int64, err error) {
	if value == nil {
		return 0, nil
	}
	vTyp := reflect.ValueOf(value)
	if vTyp.Kind() == reflect.Ptr {
		vTyp = vTyp.Elem()
		value = vTyp.Interface()
	}

	if s, ok := value.(string); ok {
		if s == "" {
			return 0, nil
		}
		value = s
	} else if b, ok := value.(json.Number); ok {
		return b.Int64()
	} else if val, ok := value.(bool); ok {
		if val {
			return 1, nil
		}
		return 0, nil
	} // 获取 val 的 reflect.Value

	// 使用 v.Kind() 来判断底层类型
	switch vTyp.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(vTyp.Int()), nil // 对于有符号整数，直接使用 Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(vTyp.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return 0, fmt.Errorf("不支持将浮点数转换为 int，请明确意图")
	default:
		res, err = convertor.ToInt(value)
		if err != nil {
			typeName := reflect.ValueOf(value).Type().String()
			return 0, errors.New("%s v is not int, %s", value, typeName)
		}
	}
	return res, err
}

func ToFloat64(value any) (float64, error) {
	v, err := convertor.ToFloat(value)
	if err != nil {
		return 0, errors.New("%v is not float64", value)
	}
	return v, nil
}

func ToBool(value any) (bool, error) {
	var val any
	if value == nil || value == "" {
		val = "false"
	} else {
		val = value
	}
	return convertor.ToBool(fmt.Sprintf("%v", val))
}

func ToDateTime(value any) (time.Time, error) {
	return timeutils.AsTime(value)
}

func ToTime(value any) (time.Time, error) {
	return timeutils.AsTime(value)
}
