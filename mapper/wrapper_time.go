package mapper

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"reflect"
	"time"
)

type TimeWrapper struct {
	BaseTypeWrapper
}

func (w *TimeWrapper) IsType(value reflect.Value) bool {
	v := value.Interface()
	switch v.(type) {
	case *time.Time:
		return true
	case time.Time:
		return true
	}
	return false
}

func (w *TimeWrapper) SetValue(fromFieldInfo reflect.Value, toFieldInfo reflect.Value) (bool, error) {
	v := fromFieldInfo.Interface()
	switch v.(type) {
	case *time.Time:
		return w.setTime(v.(*time.Time), toFieldInfo)
	case time.Time:
		timeValue := v.(time.Time)
		return w.setTime(&timeValue, toFieldInfo)
	}
	return false, nil
}

func (w *TimeWrapper) setTime(value *time.Time, toFieldInfo reflect.Value) (bool, error) {
	switch toFieldInfo.Kind() {
	case reflect.Pointer:
		switch toFieldInfo.Interface().(type) {
		case *types.JSONTime:
			if value != nil {
				vt := types.JSONTime(*value)
				toFieldInfo.Set(reflect.ValueOf(&vt))
			} else {
				var vt *types.JSONTime
				toFieldInfo.Set(reflect.ValueOf(vt))
			}
			return true, nil
		case *time.Time:
			if value != nil {
				toFieldInfo.Set(reflect.ValueOf(value))
			} else {
				var vt *time.Time
				toFieldInfo.Set(reflect.ValueOf(vt))
			}
			return true, nil
		}
	case reflect.Struct:
		switch toFieldInfo.Interface().(type) {
		case types.JSONTime:
			vt := types.JSONTime(*value)
			toFieldInfo.Set(reflect.ValueOf(vt))
			return true, nil
		case time.Time:
			toFieldInfo.Set(reflect.ValueOf(*value))
			return true, nil
		}
	}
	return false, nil
}

func NewTimeWrapper() *TimeWrapper {
	return &TimeWrapper{}
}
