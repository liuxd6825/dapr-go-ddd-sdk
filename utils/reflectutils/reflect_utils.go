package reflectutils

import (
	"errors"
	"fmt"
	errors2 "github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
)

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

// MappingStruct
// @Description: 映射结构，将源结构属性映射到目标结构
// @param source 源结构实例
// @param result 结果结构实例
// @param set 设置方法
// @return resErr 错误
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

// MappingSlice
// @Description: 映射切片，将源切片元素映射到目标元素上。
// @param sourceSlice 源切片
// @param resultSlice  目标切片
// @param setItem 设置函数
// @return resErr 返回错误
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
	/*	if sourceSliceValue.Kind() == reflect.Ptr {
			//sourceSliceValue = sourceSliceValue.Elem()
		}
		if sourceSliceValue.Kind() == reflect.Slice {
			println("sourceSlice reflect.Slice")
			//sourceSliceValue = sourceSliceValue.Elem()
		}
	*/

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

// NewObject 支持map和struct的创建
func NewObject[T any]() (res T, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	var null T
	t := reflect.TypeOf(null)
	if t.Kind() == reflect.Map {
		v := map[string]any{}
		var a any = v
		if data, ok := a.(T); ok {
			return data, nil
		} else {
			panic("can not cast any")
		}
	}
	v, err := New(t)
	if err != nil {
		return null, err
	}
	return v.Interface().(T), nil
}

func IsMap[T any]() bool {
	var null T
	t := reflect.TypeOf(null)
	if t.Kind() == reflect.Map {
		return true
	}
	return false
}

func NewStruct[T any]() (res T, resErr error) {
	defer func() {
		if err := errors2.GetError(recover()); err != nil {
			resErr = err
		}
	}()
	var null T
	t := reflect.TypeOf(null)
	v, err := New(t)
	if err != nil {
		return null, err
	}
	return v.Interface().(T), nil
}

// NewSlice
// @Description: 动态创建切片  NewSlice[[]Object]()
// @return interface{}
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

// GetValuePointer
// @Description: 检查指针层级，只保留最后的指值
// @param v
// @return reflect.Value
func GetValuePointer(data interface{}) reflect.Value {
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

func GetFieldString(data any, fieldName string) string {
	if m, ok := data.(map[string]interface{}); ok {
		v := m[fieldName]
		return fmt.Sprintf("%v", v)
	}

	refVal := reflect.ValueOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	// 根据字段类型设置值
	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() // 设置字符串字段的值
	default:
		fmt.Println("Unsupported field type:", fieldName)
	}
	return ""
}

func SetFieldString(data any, fieldName string, val string) {
	if m, ok := data.(map[string]interface{}); ok {
		m[fieldName] = val
		return
	}

	fieldName = stringutils.FirstUpper(fieldName)
	refVal := reflect.ValueOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	// 根据字段类型设置值
	switch fieldValue.Kind() {
	case reflect.String:
		if fieldValue.CanSet() {
			fieldValue.SetString(val)
		} else {
			panic("SetString cannot set String")
		}
	default:
		fmt.Println("Unsupported field type:", fieldName)
	}
}

func SetField(data any, fieldName string, val any) {
	if m, ok := data.(map[string]interface{}); ok {
		m[fieldName] = val
		return
	}

	refVal := reflect.ValueOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	if !fieldValue.CanSet() {
		panic("SetField cannot set")
	}
	fieldValue.Set(reflect.ValueOf(val))
}

func GetField(data any, fieldName string) any {
	if m, ok := data.(map[string]interface{}); ok {
		return m[fieldName]
	}

	refVal := reflect.ValueOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	if !fieldValue.CanSet() {
		panic("SetField cannot set")
	}
	return fieldValue.Interface()
}
