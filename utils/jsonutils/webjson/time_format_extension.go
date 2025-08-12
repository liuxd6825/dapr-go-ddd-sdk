package webjson

import jsoniter "github.com/json-iterator/go"

// --- 扩展 (现在能识别 time.Time 和 *time.Time) ---
type timeFormatExtension struct {
	jsoniter.DummyExtension
}

func (extension *timeFormatExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		// 获取字段的类型字符串
		field := binding.Field
		fieldType := field.Type()
		typeName := fieldType.String()
		switch typeName {
		case "time.Time", "times.Time":
			tf := &timeEncoder{}
			binding.Encoder = tf
			binding.Decoder = tf
		case "*time.Time", "*times.Time":
			ptf := &pointerTime{}
			binding.Encoder = ptf
			binding.Decoder = ptf
		case "times.Date":
			tf := &timeEncoder{}
			binding.Encoder = tf
			binding.Decoder = tf
		case "*times.Date":
			ptf := &pointerTime{}
			binding.Encoder = ptf
			binding.Decoder = ptf
		}
	}
}
