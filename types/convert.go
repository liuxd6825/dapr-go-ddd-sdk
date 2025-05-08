package types

import (
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"time"
)

type ConvertType = string

const (
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
	case ConvertTypeInt:
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
		return nil, errors.New(fmt.Sprintf("unknown convert type: %s", convType))
	}
}

func ConvertString(value any) (string, error) {
	v := convertor.ToString(value)
	return v, nil
}

func ConvertInt(value any) (int64, error) {
	v, err := convertor.ToInt(value)
	if err != nil {
		return 0, errors.New("%sv is not int", v)
	}
	return int64(v), nil
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
