package server

import (
	"encoding/json"
	"regexp"
	"strings"
)

type JsonUtils struct {
}

func NewJsonUtils() *JsonUtils {
	return &JsonUtils{}
}

func (j *JsonUtils) NewMap(jsonStr string) map[string]any {
	jsonStr = preprocessJSON(jsonStr)
	data := make(map[string]any)
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}
	return data
}

// preprocessJSON 修正不合法 JSON 的函数
func preprocessJSON(jsonStr string) string {
	// 1. 为字段名添加双引号
	reFields := regexp.MustCompile(`(\$\w+|[a-zA-Z_]\w*):`)
	jsonStr = reFields.ReplaceAllString(jsonStr, `"$1":`)

	// 2. 将单引号替换为双引号
	reQuotes := regexp.MustCompile(`'([^']*)'`)
	jsonStr = reQuotes.ReplaceAllString(jsonStr, `"$1"`)

	// 3. 删除多余的逗号
	reTrailingComma := regexp.MustCompile(`,(\s*[}\]])`)
	jsonStr = reTrailingComma.ReplaceAllString(jsonStr, `$1`)

	// 4. 删除多余的空白
	jsonStr = strings.TrimSpace(jsonStr)

	return jsonStr
}
