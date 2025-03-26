package reflectutils

import (
	"errors"
	"fmt"
	errors2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

// NewInstance 泛型函数：根据类型参数 T 创建实例
func NewInstance[T any]() T {
	var t T
	// 获取类型 T 的反射类型
	typ := reflect.TypeOf(t)

	// 如果 T 是指针类型，创建指针指向的实例
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()         // 获取指针指向的类型
		v := reflect.New(typ)    // 创建新实例
		return v.Interface().(T) // 转换为 T 类型
	}

	// 如果 T 是非指针类型，直接创建实例
	return reflect.New(typ).Elem().Interface().(T)
}

// IsMapStringKey 泛型函数：检查 T 是否是 map[string]... 类型
func IsMapStringKey[T any](t T) bool {
	// 获取类型 T 的反射类型
	typ := reflect.TypeOf(t)

	// 检查是否是 map 类型
	if typ.Kind() != reflect.Map {
		return false
	}

	// 检查 key 的类型是否是 string
	return typ.Key().Kind() == reflect.String
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

func IsStruct[T any]() bool {
	var null T
	t := reflect.TypeOf(null)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		return true
	}
	return false
}

// GetStructName 获取 any 变量的结构体类型名称（不包含包路径）
func GetStructName(v any) string {
	// 获取反射类型
	t := reflect.TypeOf(v)

	// 处理指针类型
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem() // 解引用指针
	}

	// 检查是否为结构体类型
	if t != nil && t.Kind() == reflect.Struct {
		return t.Name() // 返回完整的类型名称（包路径 + 类型名）
	}
	return ""
}

// GetStructFullName 获取 any 变量的结构体类型名称（包含包路径）
func GetStructFullName(v any) string {
	// 获取反射类型
	t := reflect.TypeOf(v)

	// 处理指针类型
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem() // 解引用指针
	}

	// 检查是否为结构体类型
	if t != nil && t.Kind() == reflect.Struct {
		return t.String() // 返回完整的类型名称（包路径 + 类型名）
	}
	return ""
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

	refVal := ValueElemOf(data)
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

func GetClassName[T any]() string {
	var null T
	t := reflect.TypeOf(null)
	// 处理指针类型
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem() // 解引用指针
	}

	// 检查是否为结构体类型
	if t.Kind() == reflect.Struct {
		return t.String() // 返回完整的类型名称（包路径 + 类型名）
	} else if t.Kind() == reflect.Map {
		return t.String()
	}
	return ""
}

func SetFieldString(data any, fieldName string, val string) {
	if m, ok := data.(map[string]interface{}); ok {
		m[fieldName] = val
		return
	}

	refVal := ValueElemOf(data)
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

func SetField(data any, fieldName string, val any) bool {
	if m, ok := data.(map[string]interface{}); ok {
		m[fieldName] = val
		return true
	}

	refVal := ValueElemOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	if fieldValue.IsValid() {
		if !fieldValue.CanSet() {
			return false
		}
		fieldValue.Set(reflect.ValueOf(val))
		return true
	}
	return true
}

func GetField(data any, fieldName string) any {
	if m, ok := data.(map[string]interface{}); ok {
		return m[fieldName]
	}

	refVal := ValueElemOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	if !fieldValue.CanSet() {
		panic("SetField cannot set")
	}
	return fieldValue.Interface()
}

func GetFieldStrings(data any, fieldName string) []string {
	if m, ok := data.(map[string]interface{}); ok {
		v, ok := m[fieldName]
		if !ok || v == nil {
			return []string{}
		}
		if list, ok := v.([]string); ok {
			return list
		} else {
			panic("GetFieldStrings cannot get []string")
		}
	}

	refVal := ValueElemOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	// 根据字段类型设置值
	switch fieldValue.Kind() {
	case reflect.Slice:
		val := fieldValue.Interface() // 设置字符串字段的值
		if val == nil {
			return []string{}
		}
		list, ok := val.([]string)
		if !ok {
			panic("GetFieldStrings cannot get []string")
		}
		if len(list) == 0 {
			return []string{}
		}
		return list
	default:
		fmt.Println("Unsupported field type:", fieldName)
	}
	return []string{}
}

func SetFieldStrings(data any, fieldName string, val []string) {
	if m, ok := data.(map[string]interface{}); ok {
		m[fieldName] = val
		return
	}

	fieldName = stringutils.FirstUpper(fieldName)
	refVal := ValueElemOf(data)
	// 根据字段名称获取字段的反射值
	fieldValue := refVal.FieldByName(fieldName)
	// 根据字段类型设置值
	switch fieldValue.Kind() {
	case reflect.Slice:
		if fieldValue.CanSet() {
			fieldValue.Set(reflect.ValueOf(val))
		} else {
			panic("SetString cannot set String")
		}
	default:
		fmt.Println("Unsupported field type:", fieldName)
	}
}

func ValueElemOf(data any) reflect.Value {
	refVal := reflect.ValueOf(data)
	for {
		if refVal.Kind() == reflect.Ptr {
			refVal = refVal.Elem()
		} else {
			break
		}
	}
	return refVal
}
