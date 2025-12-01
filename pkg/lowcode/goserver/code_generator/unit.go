package code_generator

import (
	"strings"
	"unicode"

	"github.com/liuxd6825/jsonschema/v6"
)

// ==========================================
// 3. 辅助函数 (字符串处理)
// ==========================================

func ToPascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// 简单处理：首字母大写，如果遇到 _ 或 - 则后一位大写
	var result strings.Builder
	upperNext := true
	for _, r := range s {
		if r == '_' || r == '-' {
			upperNext = true
			continue
		}
		if upperNext {
			result.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// GetGoType 将 JsonSchema 类型映射为 Golang 类型
func GetGoType(s *jsonschema.Schema) string {
	// 处理类型数组，例如 ["string", "null"]
	if s.Types == nil {
		panic("jsonscheam.types no nil")
	}
	var typeStr = ""
	var isNull bool = false
	typeList := s.Types.ToStrings()
	for _, t := range typeList {
		if t == "null" {
			isNull = true
		} else {
			typeStr = t
		}
	}

	var goType string
	switch typeStr {
	case "string":
		// 特殊处理日期

		goType = "string"
	case "integer", "int":
		goType = "int64"
	case "number":
		goType = "float64"
	case "boolean":
		goType = "bool"
	case "array":
		if s.Items != nil {
			//return "[]" + GetGoType(s.Items)
		}
		if s.Items2020 != nil {
			goType = "[]" + GetGoType(s.Items2020)
		} else {
			goType = "[]any"
		}
	case "date", "date-time":
		// 需要引入 time 包，这里为了演示简化，如果需要 time.Time 可以在模板中处理 import
		goType = "time.Time" // 或者 time.Time
	case "object":
		goType = "map[string]interface{}"
	default:
		// 如果是日期类型，Schema中有时不是标准type
		if typeStr == "date" {
			goType = "time.Time"
		} else {
			goType = "string"
		}
	}

	if isNull {
		return "*" + goType
	}
	return goType
}
