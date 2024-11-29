package hserver

import (
	"fmt"
	"regexp"
)

type Alias map[string]string

func NewAlias() Alias {
	return make(Alias)
}

func ParseAlias(input string) (Alias, error) {
	if len(input) == 0 {
		return Alias{}, nil
	}
	// 定义正则表达式匹配类似 "key:value" 的键值对
	re := regexp.MustCompile(`\s*([\w]+)\s*:\s*([\w]+)\s*`)
	matches := re.FindAllStringSubmatch(input, -1)

	// 检查匹配结果
	if matches == nil {
		return nil, fmt.Errorf("invalid input format")
	}

	// 创建结果 map
	result := make(map[string]string)

	// 遍历匹配项并填充 map
	for _, match := range matches {
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			result[key] = value
		}
	}

	return result, nil
}

func (a Alias) GetAlias(key string) string {
	if a == nil {
		return key
	}
	newKey, ok := a[key]
	if !ok {
		return key
	}
	return newKey
}
