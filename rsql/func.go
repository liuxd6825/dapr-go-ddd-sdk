package rsql

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var FuncTypes = []string{"sub"}

type FuncValue struct {
	Value   string
	Name    string
	Content string
	Args    map[string]string
	Rsql    string
}

func NewFuncValue(value string) (*FuncValue, error) {
	i := strings.Index(value, "{")
	name := value[:i]
	name = strings.TrimSpace(name)

	c := value[i+1:]
	content := c[:len(c)-1]
	args, err := parseQuery(content)
	if err != nil {
		return nil, err
	}
	rSql := args["rsql"]
	return &FuncValue{Value: value, Name: name, Content: content, Args: args, Rsql: rSql}, nil
}
func (f *FuncValue) ValueName() string { return "func" }

func (f *FuncValue) RSQL() string {
	return f.Rsql
}

func parseNonQuotedJSON(query string) (map[string]any, error) {
	// 使用正则表达式将所有不带双引号的键转换为带双引号的键
	re := regexp.MustCompile(`([a-zA-Z0-9_]+)\s*:`)
	query = re.ReplaceAllString(query, `"$1":`)

	// 为了避免错误，确保 JSON 结构能正确解析
	query = "{" + query + "}"

	// 定义一个 map 来存储结果
	var result map[string]interface{}
	// 使用标准库的 JSON 解析器
	err := json.Unmarshal([]byte(query), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 将类似 "table:orderItems;field:customerId;rsql:product=like=*book*" 转换为 map[string]string
func parseQuery(query string) (map[string]string, error) {
	// 创建一个空的 map 用来存储结果
	result := make(map[string]string)

	// 按 ; 分割字符串
	parts := strings.Split(query, ",")

	// 遍历每一部分并处理
	for _, part := range parts {
		// 按 : 分割每个键值对
		keyValue := strings.Split(part, ":")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", part)
		}
		// 将键值对添加到 map 中
		key := strings.TrimSpace(keyValue[0])
		val := strings.TrimSpace(keyValue[1])
		result[key] = val
	}

	return result, nil
}
