package webjson

import (
	"fmt"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// --- time.Time 的编码器/解码器 (无需修改) ---
type timeEncoder struct{}

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
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		stream.WriteString(t.Format(GlobalDefaultDataFormat))
		return
	}
	stream.WriteString(t.Format(GlobalDefaultTimeFormat))
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
	t, err := time.Parse(GlobalDefaultTimeFormat, val)
	if err != nil {
		iter.Error = fmt.Errorf("timeDecoder: failed to parse time string '%s' with format '%s': %w", val, GlobalDefaultTimeFormat, err)
		return
	}
	*((*time.Time)(ptr)) = t
}
