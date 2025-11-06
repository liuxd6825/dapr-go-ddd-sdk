package webjson

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

// =================================================================
//
//   核心：创建一个能够处理所有类型（包括 map 值）的扩展
//
// =================================================================

type mapExtension struct {
	jsoniter.DummyExtension
}

// DecorateEncoder 会在 jsoniter 为任何类型寻找编码器时被调用
func (e *mapExtension) DecorateEncoder(typ reflect2.Type, defaultEncoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	// 使用类型的字符串表示进行判断
	switch typ.String() {
	case "time.Time":
		return &timeEncoder{}
	case "times.Time":
		return &typesTimeEncoder{}
	case "times.Date":
		return &typesTimeEncoder{}
	}
	// 对指针类型的判断
	if typ.Kind() == reflect.Ptr {
		switch typ.Type1().Elem().String() {
		case "time.Time":
			return &pointerTime{}
			// case "times.Time": ... (如果需要，也为指针类型 times.Time 添加)
		}
	}

	// 对于所有其他类型，返回它们默认的编码器，这一点至关重要！
	return defaultEncoder
}

// DecorateDecoder 的逻辑与 DecorateEncoder 完全对应
func (e *mapExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	typeName := typ.String()
	switch typeName {
	case "time.Time":
		return &timeEncoder{}
	case "times.Time":
		return &typesTimeEncoder{}
	case "times.Date":
		return &typesTimeEncoder{}
	}
	if typ.Kind() == reflect.Ptr {
		switch typ.Type1().Elem().String() {
		case "time.Time":
			return &pointerTime{}
		}
	}
	return decoder
}
