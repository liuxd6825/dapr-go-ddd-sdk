package webjson

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"time"
	"unsafe"
)

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
