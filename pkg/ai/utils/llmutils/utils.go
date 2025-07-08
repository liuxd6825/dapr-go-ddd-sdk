package llmutils

import "strings"

func GetJsonString(content string) string {
	if strings.HasPrefix(content, "```json") && strings.HasSuffix(content, "```") {
		return content[7 : len(content)-3]
	}
	return content
}
