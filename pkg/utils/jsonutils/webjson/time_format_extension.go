package webjson

import jsoniter "github.com/json-iterator/go"

// --- 扩展 (现在能识别 time.Time 和 *time.Time) ---
type TimeFormatExtension struct {
	jsoniter.DummyExtension
	tf  *timeEncoder
	ptf *pointerTime
}

func NewTimeFormatExtension() *TimeFormatExtension {
	return &TimeFormatExtension{
		tf:  &timeEncoder{},
		ptf: &pointerTime{},
	}
}

func (e *TimeFormatExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		// 获取字段的类型字符串
		field := binding.Field
		fieldType := field.Type()
		typeName := fieldType.String()
		switch typeName {
		case "time.Time", "times.Time":
			binding.Encoder = e.tf
			binding.Decoder = e.tf
		case "*time.Time", "*times.Time":
			binding.Encoder = e.ptf
			binding.Decoder = e.ptf
		case "times.Date":
			binding.Encoder = e.tf
			binding.Decoder = e.tf
		case "*times.Date":
			binding.Encoder = e.ptf
			binding.Decoder = e.ptf
		}
	}
}
