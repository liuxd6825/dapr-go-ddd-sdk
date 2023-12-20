package jsonutils

import (
	"fmt"
	"github.com/goccy/go-json"
	"strings"
	"time"
)

// Marshal  格式化json字符串
func Marshal(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

// MarshalIndent  转化并格式化json字符串
func MarshalIndent(data interface{}) (string, error) {
	bs, err := json.MarshalIndent(data, "", "    ")
	return string(bs), err
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// MarshalNoKeyMarks key没有双引号，序列化json字符串。
func MarshalNoKeyMarks(data map[string]interface{}) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("{")
	count := len(data)
	i := 0
	for k, v := range data {
		sb.WriteString(k + ":")
		switch v.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("'%v'", v))
		case *string:
			sb.WriteString(fmt.Sprintf("'%v'", v))
		case time.Time:
			sb.WriteString(fmt.Sprintf("'%v'", v))
		case *time.Time:
			sb.WriteString(fmt.Sprintf("'%v'", v))
		case map[string]interface{}:
			props := v.(map[string]interface{})
			if jsonStr, err := MarshalNoKeyMarks(props); err != nil {
				return "", err
			} else {
				sb.WriteString(jsonStr)
			}
		case nil:
			sb.WriteString("null")
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}
		if i < count-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String(), nil
}
