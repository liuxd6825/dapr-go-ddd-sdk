package jsonutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ParseTime = func(val string, key any) (any, error)

type UnmarshalTimeOptions struct {
	ParseTime  ParseTime
	TimeFields map[string]any
}

func UnmarshalTime(data []byte, opts ...*UnmarshalTimeOptions) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var opt *UnmarshalTimeOptions
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &UnmarshalTimeOptions{}
	}
	return parseValue(decoder, opt)
}

func parseValue(decoder *json.Decoder, opts *UnmarshalTimeOptions) (any, error) {
	// 获取下一个 token
	t, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t := t.(type) {
	case json.Delim:
		// 如果是对象
		if t == '{' {
			return parseObject(decoder, opts)
		}
		// 如果是数组
		if t == '[' {
			return parseArray(decoder, opts)
		}
	}

	// 如果是基础值，直接返回
	return t, nil
}

func parseObject(decoder *json.Decoder, opts *UnmarshalTimeOptions) (map[string]any, error) {
	obj := make(map[string]any)

	for decoder.More() {
		// 读取 key
		t, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		key, ok := t.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for key, got %T", t)
		}

		// 获取子字段的 TimeFields 配置
		var subFields map[string]any
		var hasSubFields = false
		var subVal any
		if opts.TimeFields != nil {
			subVal, hasSubFields = opts.TimeFields[key]
			if hasSubFields {
				if sub, ok := subVal.(map[string]any); ok {
					subFields = sub
				}
			}
		}

		// 读取 value
		value, err := parseValue(decoder, &UnmarshalTimeOptions{ParseTime: opts.ParseTime, TimeFields: subFields})
		if err != nil {
			return nil, err
		}

		// 如果 key 是需要解析为时间的字段
		if hasSubFields {
			if strVal, ok := value.(string); ok {
				if opts.ParseTime != nil {
					value, err = opts.ParseTime(strVal, subVal)
					if err != nil {
						return nil, err
					}
				} else {
					const dateFormat = "2006-01-02 15:04:05"
					if parsedTime, err := time.Parse(dateFormat, strVal); err == nil {
						value = parsedTime
					} else {
						return nil, fmt.Errorf("failed to parse time for field %s: %v", key, err)
					}
				}
			}
		}
		if n, ok := value.(json.Number); ok {
			if strings.Contains(n.String(), ".") {
				value, _ = n.Float64()
			} else {
				value, _ = n.Int64()
			}
		}
		obj[key] = value
	}

	// 读取结束标志 '}'
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}

	return obj, nil
}

func parseArray(decoder *json.Decoder, opts *UnmarshalTimeOptions) ([]any, error) {
	var arr []any

	for decoder.More() {
		// 解析数组中的每个元素
		value, err := parseValue(decoder, opts)
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}

	// 读取结束标志 ']'
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}

	return arr, nil
}
