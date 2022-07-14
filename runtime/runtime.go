package runtime

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

	if t.Kind() == reflect.Slice {
		e := t.Elem()
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}
		return e
	}
	return nil
}

func MappingSlice(sourceSlice interface{}, targetSlice interface{}, setItem func(index int, source reflect.Value, target reflect.Value) error) (resErr error) {
	defer func() {
		if err := errorutils.GetError(recover()); err != nil {
			resErr = err
		}
	}()

	if setItem == nil {
		return fmt.Errorf("MappingSlice(sourceSlice, targetSlice, setItem) setItem is nil")
	}

	itemType := NewSliceItemType(targetSlice)
	if itemType == nil {
		return fmt.Errorf("MappingSlice(sourceSlice, targetSlice, setItem) itemType is nil")
	}

	targetSliceValue := reflect.ValueOf(targetSlice)
	if targetSliceValue.Kind() == reflect.Ptr {
		targetSliceValue = targetSliceValue.Elem()
	}
	if targetSliceValue.Kind() == reflect.Slice {
		println("targetSliceValue.Kind is reflect.Slice")
		//sourceSliceValue = sourceSliceValue.Elem()
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
		target := reflect.New(itemType)
		source := sourceSliceValue.Index(i)
		if err := setItem(i, source, target); err != nil {
			return fmt.Errorf("MappingSlice(sourceSlice, targetSlice, setItem) set index %v error: %s", i, err.Error())
		}
		if itemType.Kind() == reflect.Ptr {
			targetSliceValue = reflect.Append(targetSliceValue, target)
		} else {
			targetSliceValue = reflect.Append(targetSliceValue, target.Elem())
		}
	}
	return nil
}
