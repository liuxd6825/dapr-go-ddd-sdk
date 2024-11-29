package runtime

import (
	"reflect"
	"strings"
)

type FieldNameMapper struct{}

// FieldName 映射字段名称
func (m *FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return m.lowerFirstLetter(f.Name) // 将字段名转换为小写
}

// MethodName 映射方法名称
func (m *FieldNameMapper) MethodName(t reflect.Type, mtd reflect.Method) string {
	return m.lowerFirstLetter(mtd.Name) // 将方法名转换为小写
}

// 辅助函数：将首字母转换为小写
func (m *FieldNameMapper) lowerFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
