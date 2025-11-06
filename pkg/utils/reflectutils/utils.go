package reflectutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

func ParseDescMap(data any) (map[string]string, error) {
	res := map[string]string{}
	obj := NewRefObj(data)
	fields := obj.FieldsAll()
	for _, f := range fields {
		desc, err := f.Tag("desc")
		if err != nil {
			return nil, err
		}
		if desc == "" {
			desc = f.Name()
		}
		val, err := f.Get()
		if err != nil {
			return nil, err
		}
		res[desc] = stringutils.AnyToString(val)
	}
	return res, nil
}
func ParseDesc(data any) (string, error) {
	sb, err := ParseDescOptions(data, "desc", "ulog")
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func ParseDescOptions(data any, nameTag string, ulogTag string) (*strings.Builder, error) {
	res := &strings.Builder{}
	obj := NewRefObj(data)
	fields := obj.FieldsAll()

	for _, f := range fields {
		userLog, err := f.Tag(ulogTag)
		if err != nil {
			return nil, err
		}
		if userLog == "-" {
			continue
		}

		desc, err := f.Tag(nameTag)
		if err != nil {
			return nil, err
		}
		if desc == "" {
			desc = f.Name()
		}
		val, err := f.Get()
		if err != nil {
			return nil, err
		}

		res.WriteString(fmt.Sprintf("[%s]=`%s`; ", desc, stringutils.AnyToString(val)))
	}

	return res, nil
}

// ParseTag parses a golang struct tag into a map.
func ParseTag(tag string) (map[string]string, error) {
	res := map[string]string{}

	// This code is copied/modified from: reflect/type.go:
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, fmt.Errorf("Cannot unquote tag %s in %s: %s", name, tag, err.Error())
		}
		res[name] = value
	}

	return res, nil
}

// GetListFieldValues 从 any 类型的列表中提取指定字段的值，并返回 []string
func GetListFieldValues(list any, field string) ([]string, error) {
	// 使用反射获取 data 的类型
	v := reflect.ValueOf(list)

	// 判断 data 是否是切片或数组
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("input is not a slice or array")
	}

	// 准备结果切片
	result := make([]string, 0, v.Len())

	// 遍历切片或数组中的每个元素
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		// 判断元素是结构体还是 map
		switch elem.Kind() {
		case reflect.Struct:
			// 如果是结构体，获取字段值
			fieldValue := elem.FieldByName(field)
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("field %s not found in struct", field)
			}
			result = append(result, fmt.Sprintf("%v", fieldValue.Interface()))
		case reflect.Map:
			// 如果是 map，获取键对应的值
			fieldValue := elem.MapIndex(reflect.ValueOf(field))
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("field %s not found in map", field)
			}
			result = append(result, fmt.Sprintf("%v", fieldValue.Interface()))
		default:
			return nil, fmt.Errorf("unsupported element type: %s", elem.Kind())
		}
	}

	return result, nil
}
