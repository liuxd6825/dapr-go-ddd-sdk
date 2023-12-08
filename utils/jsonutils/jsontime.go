package jsonutils

import (
	"github.com/json-iterator/go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"time"
	"unsafe"
)

const (
	fDatetime = "datetime"
	fDate     = "date"
)

var CustomJson = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	CustomJson.RegisterExtension(&CustomTimeExtension{})
}

type CustomTimeExtension struct {
	jsoniter.DummyExtension
}

func (e *CustomTimeExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {

	for _, binding := range structDescriptor.Fields {
		var typeErr error
		var isPtr bool
		typeName := binding.Field.Type().String()

		if typeName == "time.Time" {
			isPtr = false
		} else if typeName == "*time.Time" {
			isPtr = true
		} else {
			continue
		}

		formatTag := binding.Field.Tag().Get("time_format")
		if len(formatTag) == 0 {
			formatTag = fDatetime
		}
		timeFormat := formatTag
		if timeFormat == fDatetime {
			timeFormat = "2006-01-02 15:04:05"
		} else if timeFormat == fDate {
			timeFormat = "2006-01-02"
		}

		locale := setting.GetTimeZone()
		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if typeErr != nil {
				stream.Error = typeErr
				return
			}

			var format string
			if timeFormat == "" {
				// format = time.RFC3339Nano
				format = "2006-01-02 15:04:05"
			} else {
				format = timeFormat
			}

			var tp *time.Time
			if isPtr {
				tpp := (**time.Time)(ptr)
				tp = *(tpp)
			} else {
				tp = (*time.Time)(ptr)
			}

			if tp != nil {
				lt := tp.In(locale)
				str := lt.Format(format)
				if formatTag == fDate && (str == "0000-01-01" || (lt.Unix() <= 0)) {
					str = "0000-00-00"
				} else if formatTag == fDatetime && (str == "0000-01-01 00:00:00" || (lt.Unix() <= 0)) {
					str = "0000-00-00 00:00:00"
				}
				stream.WriteString(str)
			} else {
				_, _ = stream.Write([]byte("null"))
			}
		}}

		binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			if typeErr != nil {
				iter.Error = typeErr
				return
			}

			var format string
			if timeFormat == "" {
				format = time.RFC3339
			} else {
				format = timeFormat
			}

			str := iter.ReadString()
			var t *time.Time
			if str != "" {
				var err error
				tmp, err := time.ParseInLocation(format, str, locale)
				if err != nil {
					if _, ok := err.(*time.ParseError); ok {
						tmp = time.Date(0, 1, 1, 0, 0, 0, 0, locale)
					} else {
						iter.Error = err
						return
					}
				}
				t = &tmp
			} else {
				t = nil
			}

			if isPtr {
				tpp := (**time.Time)(ptr)
				*tpp = t
			} else {
				tp := (*time.Time)(ptr)
				if tp != nil && t != nil {
					*tp = *t
				}
			}
		}}
	}
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}
