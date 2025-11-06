package webjson

import (
	jsoniter "github.com/json-iterator/go"
)

// GlobalDefaultTimeFormat 定义全局默认日期格式 (无需修改)
const GlobalDefaultTimeFormat = "2006-01-02 15:04:05"
const GlobalDefaultDataFormat = "2006-01-02"

// JSON 是我们定制的、具有全局默认时间格式的json-iterator实例
var JSON jsoniter.API

func init() {

	// 我们还需要为指针类型和自定义类型手动注册编码器，以确保万无一失。
	// jsoniter 在处理 interface{} 时，会依赖这些全局注册。
	/*
		tEncoder := &timeEncoder{}
		jsoniter.RegisterTypeEncoder("time.Time", tEncoder)
		//jsoniter.RegisterTypeDecoder("time.Time", tEncoder)

		ptEncoder := &pointerTime{}
		jsoniter.RegisterTypeEncoder("*time.Time", ptEncoder)
		jsoniter.RegisterTypeDecoder("time.Time", ptEncoder)

		typeEncoder := &typesTimeEncoder{}
			jsoniter.RegisterTypeEncoder("times.Time", typeEncoder)
			jsoniter.RegisterTypeDecoder("times.Time", typeEncoder)

			jsoniter.RegisterTypeEncoder("times.Date", typeEncoder)
			jsoniter.RegisterTypeDecoder("times.Date", typeEncoder)
	*/

	// JSON 是我们导出的预先配置好的 JSON API 实例
	JSON = jsoniter.Config{
		// 1. 手动设置与标准库兼容的标志
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		// 3. 最后，调用 .Froze() "冻结"整个配置，生成最终的API
	}.Froze()
	// 在此基础上，为 time.Time 类型注册我们自定义的“默认”编码器
	// 这个编码器会在没有 time_format 标签时被调用
	JSON.RegisterExtension(NewTimeFormatExtension())
	JSON.RegisterExtension(&mapExtension{})

}

func Marshal(data interface{}) (string, error) {
	return JSON.MarshalToString(data)
}

func Unmarshal(data []byte, v any) error {
	return JSON.Unmarshal(data, v)
}
