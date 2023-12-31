package reflectutils

import (
	"errors"
	"fmt"
	errors2 "github.com/liuxd6825/dapr-go-ddd-sdk/errors"
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

//
// NewSliceItemType
// @Description: 根据给定切片，返回元素的Type
// @param slice
// @return reflect.Type
func NewSliceItemType(slice interface{}) (res reflect.Type, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	if slice == nil {
		return nil, errors.New("slice is nil")
	}
	t := reflect.TypeOf(slice)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		return nil, errors.New("slice kind not is reflect.Slice")
	}
	e := t.Elem()
	return e, nil
}

//
// MappingStruct
// @Description: 映射结构，将源结构属性映射到目标结构
// @param source 源结构实例
// @param result 结果结构实例
// @param set 设置方法
// @return resErr 错误
//
func MappingStruct(source interface{}, result interface{}, set func(source, target reflect.Value) error) (resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	if set == nil {
		return fmt.Errorf("MappingStruct(sourceSlice, targetSlice, setItem) setItem is nil")
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
		return fmt.Errorf("MappingStruct(source, result, set) error: %s", err.Error())
	}
	resultsValue.Elem().Set(targetValue)

	return nil
}

//
// MappingSlice
// @Description: 映射切片，将源切片元素映射到目标元素上。
// @param sourceSlice 源切片
// @param resultSlice  目标切片
// @param setItem 设置函数
// @return resErr 返回错误
//
func MappingSlice(sourceSlice interface{}, resultSlice interface{}, setItem func(index int, source reflect.Value, target reflect.Value) error) (resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()

	if setItem == nil {
		return fmt.Errorf("MappingSlice(sourceSlice, resultSlice, setItem) setItem is nil")
	}

	resultsValue := reflect.ValueOf(resultSlice)
	if resultsValue.Kind() != reflect.Ptr {
		return fmt.Errorf("MappingSlice(sourceSlice, resultSlice, setItem) err: resultSlice must be a pointer to a slice, but was a %s", resultsValue.Kind())
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

func New(t reflect.Type) (res reflect.Value, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	if t.Kind() == reflect.Ptr {
		elem := reflect.New(t.Elem())
		return elem, nil
	}
	return reflect.New(t), nil
}

func NewStruct[T interface{}]() (res T, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	var null T
	v, err := New(reflect.TypeOf(null))
	if err != nil {
		return null, err
	}
	return v.Interface().(T), nil
}

//
// NewSlice
// @Description: 动态创建切片  NewSlice[[]Object]()
// @return interface{}
//
func NewSlice[T interface{}]() (res T, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	var t T
	v, err := New(reflect.TypeOf(t))
	if err != nil {
		return t, err
	}

	res, _ = v.Elem().Interface().(T)
	return res, nil
}

//
// GetValuePointer
// @Description: 检查指针层级，只保留最后的指值
// @param v
// @return reflect.Value
//
func GetValuePointer(data interface{}) reflect.Value {
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}
