package json_pkg

import (
	"bytes"
	"encoding/json"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/jsonutils"
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
func (j *JsonPkg) Format(rawJSON any) []byte {
	data := pkg.ToBytes(rawJSON)
	// 创建一个缓冲区来存储格式化后的 JSON
	var formattedJSON bytes.Buffer
	// 使用 json.Indent 进行格式化
	err := json.Indent(&formattedJSON, data, "", "  ") // 第二个参数是缩进字符串
	if err != nil {
		panic(errors.New("json format error: %s", err.Error()))
	}
	return formattedJSON.Bytes()
}

// NewMap
//
//	@Description:
//	@receiver j
//	@param jsonVal 可以是string或[]byte类型
//	@return any
func (j *JsonPkg) NewMap(jsonVal any) any {
	var res any = make(map[string]any)
	var data []byte
	if str, ok := jsonVal.(string); ok {
		data = []byte(str)
	} else if b, ok := jsonVal.([]byte); ok {
		data = b
	} else {
		panic(errors.New("jsonVal is not string"))
	}

	err := jsonutils.Unmarshal(data, &res)
	if err != nil {
		panic(errors.New("json parse error: %s", err.Error()))
	}
	return res
}
