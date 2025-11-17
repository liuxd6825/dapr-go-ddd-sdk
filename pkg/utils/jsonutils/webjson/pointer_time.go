package webjson

import (
	"fmt"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

// !!! 新增部分：*time.Time 的编码器和解码器 !!!
type pointerTime struct {
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
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		stream.WriteString(t.Format(GlobalDefaultDataFormat))
		return
	}
	stream.WriteString(t.Format(GlobalDefaultTimeFormat))
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
	if len(val) <= 10 {
		t, err := timeutils.AsTime(val)
		if err != nil {
			iter.Error = fmt.Errorf("pointerTimeDecoder: failed to parse time string '%s' with format '%s': %w", val, GlobalDefaultDataFormat, err)
			return
		}
		// 将解析出的 time.Time 的地址赋给指针
		*(**time.Time)(ptr) = &t
		return
	}
	t, err := timeutils.AsTime(val)
	if err != nil {
		iter.Error = fmt.Errorf("pointerTimeDecoder: failed to parse time string '%s' with format '%s': %w", val, GlobalDefaultTimeFormat, err)

	}
	// 将解析出的 time.Time 的地址赋给指针
	*(**time.Time)(ptr) = &t
}
