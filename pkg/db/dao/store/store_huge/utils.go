package store_huge

import (
	"fmt"
	"strings"
)

// formatScript 将 bindings 嵌入 Gremlin 脚本,避免依赖 SDK 中尚未导出的 WithBindings 选项
// bindings 必须是字面量安全的值(string/number/bool);否则请预先 json.Marshal
func formatScript(script string, bindings map[string]any) string {
	if len(bindings) == 0 {
		return script
	}
	for k, v := range bindings {
		placeholder := "$<" + k + ">"
		if !strings.Contains(script, placeholder) {
			continue
		}
		script = strings.ReplaceAll(script, placeholder, literalOf(v))
	}
	return script
}

func literalOf(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		s := strings.ReplaceAll(x, "\\", "\\\\")
		s = strings.ReplaceAll(s, "'", "\\'")
		return "'" + s + "'"
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", x)
	}
}
