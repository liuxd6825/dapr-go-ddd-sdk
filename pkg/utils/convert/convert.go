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
	}

	res, err = convertor.ToInt(value)
	if err != nil {
		typeName := reflect.ValueOf(value).Type().String()
		return 0, errors.New("%s v is not int, %s", value, typeName)
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
