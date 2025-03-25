package convert

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"reflect"
	"time"
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
		return ConvertInt(value)
	case ConvertTypeFloat:
		return ConvertFloat(value)
	case ConvertTypeString,
		ConvertTypeNone:
		return ConvertString(value)
	case ConvertTypeBool:
		return ConvertBool(value)
	case ConvertTypeDateTime:
		return ConvertDateTime(value)
	case ConvertTypeTime:
		return ConvertTime(value)
	default:
		typeName := reflect.ValueOf(value).Type().String()
		return nil, errors.New(fmt.Sprintf("unknown convert type: %s , %s", convType, typeName))
	}
}

func ConvertString(value any) (string, error) {
	v := convertor.ToString(value)
	return v, nil
}

func ConvertInt(value any) (res int64, err error) {
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
	}

	res, err = convertor.ToInt(value)
	if err != nil {
		typeName := reflect.ValueOf(value).Type().String()
		return 0, errors.New("%s v is not int, %s", value, typeName)
	}
	return res, err
}

func ConvertFloat(value any) (float64, error) {
	v, err := convertor.ToFloat(value)
	if err != nil {
		return 0, errors.New("%v is not float64", value)
	}
	return v, nil
}

func ConvertBool(value any) (bool, error) {
	return convertor.ToBool(fmt.Sprintf("%v", value))
}

func ConvertDateTime(value any) (time.Time, error) {
	return timeutils.AsTime(value)
}

func ConvertTime(value any) (time.Time, error) {
	return timeutils.AsTime(value)
}
