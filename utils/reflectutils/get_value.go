package reflectutils

import (
	"reflect"
	"strconv"
)

func ValueToString(val reflect.Value) (string, error) {
	switch val.Type().Kind() {
	case reflect.Invalid:
		break
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uintptr:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Float32:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.Complex64:
		break
	case reflect.Complex128:
		break
	case reflect.Array:
		break
	case reflect.Chan:
		break
	case reflect.Func:
		break
	case reflect.Interface:
		break
	case reflect.Map:
		break
	case reflect.Pointer:
		break
	case reflect.Slice:
		break
	case reflect.String:
		break
	case reflect.Struct:
		break
	case reflect.UnsafePointer:
		break
	}
	return "", nil
}
