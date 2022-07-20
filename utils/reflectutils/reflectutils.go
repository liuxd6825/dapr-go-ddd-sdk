package reflectutils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/errorutils"
	"reflect"
	"runtime"
)

//
// RunFuncName
// @Description: 获取当前运行的方法名称
// @param skip  调用方法的调转级数
// @return string 返回方法名称
//
func RunFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func NewSliceItemType(slice interface{}) reflect.Type {
	if slice == nil {
		return nil
	}
	t := reflect.TypeOf(slice)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Interface {
		println("t.Kind() == reflect.Interface ")
	}

	if t.Kind() == reflect.Slice {
		e := t.Elem()
		/*		if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}*/
		return e
	}
	return nil
}

func MappingOne(source interface{}, result interface{}, set func(source, target reflect.Value) error) (resErr error) {
	defer func() {
		if err := errorutils.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	if set == nil {
		return fmt.Errorf("MappingOne(sourceSlice, targetSlice, setItem) setItem is nil")
	}
	resultsValue := reflect.ValueOf(result)
	if resultsValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result argument must be a pointer to a slice, but was a %s", resultsValue.Kind())
	}

	targetValue := resultsValue.Elem()
	if targetValue.Kind() == reflect.Interface {
		targetValue = targetValue.Elem()
	}

	if targetValue.Kind() != reflect.Struct {
		return fmt.Errorf("results argument must be a pointer to a struct, but was a pointer to %s", targetValue.Kind())
	}

	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("results argument must be a pointer to a struct, but was a pointer to %s", targetValue.Kind())
	}
	if err := set(sourceValue, resultsValue); err != nil {
		return fmt.Errorf("MappingOne(source, result, set) error: %s", err.Error())
	}
	resultsValue.Elem().Set(targetValue)

	return nil
}

func MappingSlice(sourceSlice interface{}, resultSlice interface{}, setItem func(index int, source reflect.Value, target reflect.Value) error) (resErr error) {
	defer func() {
		if err := errorutils.GetError(recover()); err != nil {
			resErr = err
		}
	}()

	if setItem == nil {
		return fmt.Errorf("MappingSlice(sourceSlice, targetSlice, setItem) setItem is nil")
	}

	resultsValue := reflect.ValueOf(resultSlice)
	if resultsValue.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsValue.Kind())
	}

	sliceVal := resultsValue.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	sourceSliceValue := reflect.ValueOf(sourceSlice)
	if sourceSliceValue.Kind() == reflect.Ptr {
		//sourceSliceValue = sourceSliceValue.Elem()
	}
	if sourceSliceValue.Kind() == reflect.Slice {
		println("sourceSlice reflect.Slice")
		//sourceSliceValue = sourceSliceValue.Elem()
	}

	count := sourceSliceValue.Len()
	for i := 0; i < count; i++ {
		elementType := sliceVal.Type().Elem()
		target := reflect.New(elementType)
		source := sourceSliceValue.Index(i)
		if err := setItem(i, source, target); err != nil {
			return fmt.Errorf("MappingSlice(sourceSlice, targetSlice, setItem) set index %v error: %s", i, err.Error())
		}
		if target.Kind() == reflect.Ptr {
			sliceVal = reflect.Append(sliceVal, target.Elem())
		} else {
			sliceVal = reflect.Append(sliceVal, target)
		}
	}
	resultsValue.Elem().Set(sliceVal.Slice(0, count))
	return nil
}

func New(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem())
	}
	return reflect.New(t)
}
