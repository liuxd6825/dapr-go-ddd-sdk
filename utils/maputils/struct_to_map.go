package maputils

import (
	"fmt"
	"reflect"
	"strings"
)

// getFieldName 从 struct 字段的 tag 中解析出用作 map key 的名称
func getFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	// 如果 tag 为 "-"，表示此字段应被忽略
	if tag == "-" {
		return "", false
	}
	// `json:"name,omitempty"` -> "name"
	fieldName := strings.Split(tag, ",")[0]
	if fieldName == "" {
		// 如果没有 json tag，则使用字段名
		return field.Name, true
	}
	return fieldName, true
}

// convertValue 是核心的递归转换函数
func convertValue(v reflect.Value) any {
	// 首先处理无效值
	if !v.IsValid() {
		return nil
	}

	// 根据值的类型进行分发处理
	switch v.Kind() {
	case reflect.Ptr:
		// 如果是指针，解引用后再处理。如果指针为 nil，返回 nil。
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem())

	case reflect.Struct:
		// 如果是结构体，转换为 map[string]any
		out := make(map[string]any)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			fieldVal := v.Field(i)
			fieldTyp := t.Field(i)

			// 如果是匿名嵌入字段，则递归展开并合并
			if fieldTyp.Anonymous {
				// 确保嵌入的是结构体
				if fieldVal.Kind() == reflect.Struct {
					// 递归调用 convertValue 会返回一个 map
					embeddedMap := convertValue(fieldVal).(map[string]any)
					for key, val := range embeddedMap {
						out[key] = val
					}
				}
			} else {
				// 对于普通字段
				fieldName, ok := getFieldName(fieldTyp)
				if !ok {
					continue // 跳过 tag 为 "-" 的字段
				}
				// 递归转换字段的值
				out[fieldName] = convertValue(fieldVal)
			}
		}
		return out

	case reflect.Slice, reflect.Array:
		// 如果是切片或数组，转换为 []any
		// 如果切片为 nil，返回 nil 而不是空的[]any
		if v.IsNil() {
			return nil
		}

		out := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			// 递归转换切片中的每个元素
			out[i] = convertValue(v.Index(i))
		}
		return out

	case reflect.Interface:
		// 如果是接口，处理其内部包含的动态值
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem())

	default:
		// 对于 string, int, bool, float 等基本类型，直接返回值
		return v.Interface()
	}
}

// StructToMap 是暴露给外部的入口函数
func StructToMap(data any) (map[string]any, error) {
	val := reflect.ValueOf(data)

	// 处理指针输入
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct, but got %T", data)
	}

	result := convertValue(val)
	if m, ok := result.(map[string]any); ok {
		return m, nil
	}

	return nil, fmt.Errorf("conversion failed, result is not a map[string]any")
}
