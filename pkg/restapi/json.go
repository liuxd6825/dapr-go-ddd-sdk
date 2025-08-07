package restapi

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"time"
	"unsafe"
)

// GlobalDefaultTimeFormat 定义全局默认日期格式 (无需修改)
const GlobalDefaultTimeFormat = "2006-01-02 15:04:05"
const GlobalDefaultDataFormat = "2006-01-02"

// --- time.Time 的编码器/解码器 (无需修改) ---
type timeEncoder struct{ format string }

func (tf *timeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return (*((*time.Time)(ptr))).IsZero()
}

// (为了简洁，省略了之前已有的 Encode 和 Decode 的完整实现代码)
func (tf *timeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *((*time.Time)(ptr))
	if t.IsZero() {
		stream.WriteNil()
		return
	}
	stream.WriteString(t.Format(tf.format))
}

func (tf *timeEncoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		*((*time.Time)(ptr)) = time.Time{}
		return
	}
	val := iter.ReadString()
	if val == "" {
		*((*time.Time)(ptr)) = time.Time{}
		return
	}
	t, err := time.Parse(tf.format, val)
	if err != nil {
		iter.Error = fmt.Errorf("timeDecoder: failed to parse time string '%s' with format '%s': %w", val, tf.format, err)
		return
	}
	*((*time.Time)(ptr)) = t
}

// --- time.Time 的编码器/解码器 (无需修改) ---
type typesTime struct{ format string }

func (tf *typesTime) IsEmpty(ptr unsafe.Pointer) bool {
	return (*((*time.Time)(ptr))).IsZero()
}

// (为了简洁，省略了之前已有的 Encode 和 Decode 的完整实现代码)
func (tf *typesTime) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *((*time.Time)(ptr))
	if t.IsZero() {
		stream.WriteNil()
		return
	}
	stream.WriteString(t.Format(tf.format))
}

func (tf *typesTime) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		*((*times.Time)(ptr)) = times.Time{}
		return
	}
	val := iter.ReadString()
	if val == "" {
		*((*times.Time)(ptr)) = times.Time{}
		return
	}
	t, err := time.Parse(tf.format, val)
	if err != nil {
		iter.Error = fmt.Errorf("timeDecoder: failed to parse time string '%s' with format '%s': %w", val, tf.format, err)
		return
	}
	tVal := times.GetTime(&t)
	*((*times.Time)(ptr)) = *tVal
}

// =================================================================
//
//	!!! 新增部分：*time.Time 的编码器和解码器 !!!
//
// =================================================================
type pointerTime struct {
	format string
}

func (pte *pointerTime) IsEmpty(ptr unsafe.Pointer) bool {
	// 获取 *time.Time 指针
	t := *(**time.Time)(ptr)
	// 如果指针是 nil，或者指向一个零值，都视为空
	return t == nil || t.IsZero()
}

func (pte *pointerTime) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *(**time.Time)(ptr)
	if t == nil || t.IsZero() {
		stream.WriteNil()
		return
	}
	stream.WriteString(t.Format(pte.format))
}
func (pte *pointerTime) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		// 如果是 null, 将 Go 的 *time.Time 指针设置为 nil
		*(**time.Time)(ptr) = nil
		return
	}
	val := iter.ReadString()
	if val == "" {
		*(**time.Time)(ptr) = nil
		return
	}
	t, err := time.Parse(pte.format, val)
	if err != nil {
		iter.Error = fmt.Errorf("pointerTimeDecoder: failed to parse time string '%s' with format '%s': %w", val, pte.format, err)
		return
	}
	// 将解析出的 time.Time 的地址赋给指针
	*(**time.Time)(ptr) = &t
}

// --- 扩展 (现在能识别 time.Time 和 *time.Time) ---
type timeFormatExtension struct {
	jsoniter.DummyExtension
}

func (extension *timeFormatExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		// 获取字段的类型字符串
		typeName := binding.Field.Type().String()
		switch typeName {
		case "time.Time", "times.Time":
			tf := &timeEncoder{format: GlobalDefaultTimeFormat}
			binding.Encoder = tf
			binding.Decoder = tf
		case "*time.Time", "*times.Time":
			ptf := &pointerTime{format: GlobalDefaultTimeFormat}
			binding.Encoder = ptf
			binding.Decoder = ptf
		case "times.Date":
			tf := &timeEncoder{format: GlobalDefaultDataFormat}
			binding.Encoder = tf
			binding.Decoder = tf
		case "*times.Date":
			ptf := &pointerTime{format: GlobalDefaultDataFormat}
			binding.Encoder = ptf
			binding.Decoder = ptf
		}
	}
}

// JSON 是我们定制的、具有全局默认时间格式的json-iterator实例
var JSON jsoniter.API

func init() {

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
	JSON.RegisterExtension(&timeFormatExtension{})

}

func JsonMarshal(data interface{}) (string, error) {
	return JSON.MarshalToString(data)
}

func JsonUnmarshal(data []byte, v any) error {
	return JSON.Unmarshal(data, v)
}
