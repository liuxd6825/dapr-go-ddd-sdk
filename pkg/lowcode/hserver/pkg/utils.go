package pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/jsonutils"
)

// ToBytes
//
//	@Description: 将any转换成byte数组，支持string, []byte, map[string]any类型的转换
//	@param data
//	@return []byte
func ToBytes(data any) []byte {
	var bytes []byte = nil
	if str, ok := data.(string); ok {
		bytes = []byte(str)
	} else if bs, ok := data.([]byte); ok {
		bytes = bs
	} else if mapData, ok := data.(map[string]any); ok {
		str, err := jsonutils.Marshal(mapData)
		if err != nil {
			panic(err)
		}
		bytes = []byte(str)
	} else {
		panic("WriteFile() invalid runValues is string or []byte or map[string]any")
	}
	return bytes
}
