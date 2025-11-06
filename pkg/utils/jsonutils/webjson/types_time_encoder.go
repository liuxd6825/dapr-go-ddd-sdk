package webjson

import (
	"fmt"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
)

// --- types.Time 和 types.Date 的编码器/解码器 ---
// 假设 types.Time 和 types.Date 内部都是基于 time.Time 实现
type typesTimeEncoder struct {
}

func (tte *typesTimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return (*((*times.Time)(ptr))).Time().IsZero()
}

func (tte *typesTimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *((*times.Time)(ptr))
	tVal := t.Time()
	if tVal.IsZero() {
		stream.WriteNil()
		return
	}
	stream.WriteString(tVal.Format(GlobalDefaultTimeFormat))
}

func (tte *typesTimeEncoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if iter.ReadNil() {
		*((*times.Time)(ptr)) = times.Time{}
		return
	}
	val := iter.ReadString()
	if val == "" {
		*((*times.Time)(ptr)) = times.Time{}
		return
	}
	t, err := time.Parse(GlobalDefaultTimeFormat, val)
	if err != nil {
		iter.Error = fmt.Errorf("typesTimeDecoder: failed to parse time string '%s' with format '%s': %w", val, GlobalDefaultTimeFormat, err)
		return
	}
	// 假设有这样的转换函数
	tVal := times.GetTime(&t)
	*((*times.Time)(ptr)) = *tVal
}
