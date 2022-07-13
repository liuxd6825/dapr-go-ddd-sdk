package runtime

import (
	"fmt"
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

func SetSlice(sourceSlice interface{}, targetSlice interface{}, setItem func(source interface{}, target interface{})) (resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()

	if setItem == nil {
		return fmt.Errorf("SetSlice(targetSlice, setItem) setItem is nil")
	}

	itemType := NewSliceItemType(targetSlice)
	if itemType == nil {
		return fmt.Errorf("SetSlice(targetSlice, setItem) itemType is nil")
	}

	targetSliceValue := reflect.ValueOf(targetSlice)
	if targetSliceValue.Kind() == reflect.Ptr {
		//targetSliceValue = targetSliceValue.Elem()
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
		itemValue := reflect.New(itemType)
		source := sourceSliceValue.Index(i)
		setItem(source.Interface(), itemValue.Interface())
		if itemType.Kind() == reflect.Ptr {
			targetSliceValue = reflect.Append(targetSliceValue, itemValue)
		} else {
			targetSliceValue = reflect.Append(targetSliceValue, itemValue.Elem())
		}
	}
	return nil
}
