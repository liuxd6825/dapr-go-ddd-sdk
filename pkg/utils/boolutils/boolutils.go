package boolutils

import (
	"reflect"
	"strings"
)

func PBool(b bool) *bool {
	return &b
}

func GetBool(value any, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		return i > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		return i > 0
	case reflect.Float32, reflect.Float64:
		f := val.Float()
		return f > 0
	case reflect.String:
		s := strings.ToLower(val.String())
		return s == "true"
	default:
		panic("unhandled default case")
	}
	return defaultValue
}
