package json_pkg

import (
	"bytes"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type JsonPkg struct {
}

func New() *JsonPkg {
	return NewJsonPkg()
}

func NewJsonPkg() *JsonPkg {
	fsm := &JsonPkg{}
	return fsm
}

// Format
//
//	@Description: 对json字符串进行格式化
//	@receiver j
//	@param rawJSON
//	@return []byte
func (j *JsonPkg) Format(rawJSON []byte) []byte {
	// 创建一个缓冲区来存储格式化后的 JSON
	var formattedJSON bytes.Buffer
	// 使用 json.Indent 进行格式化
	err := json.Indent(&formattedJSON, rawJSON, "", "  ") // 第二个参数是缩进字符串
	if err != nil {
		panic(errors.New("json format error: %s", err.Error()))
	}
	return formattedJSON.Bytes()
}
