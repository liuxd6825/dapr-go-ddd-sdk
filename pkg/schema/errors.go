package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/liuxd6825/jsonschema/v6/kind"
	"strings"
)

func NewSchemaError(vError *jsonschema.SchemaValidationError) error {
	return errors.ErrorOf("schema error: %s", vError.Error())
}

func NewFieldsError(vError *jsonschema.ValidationError) error {
	err := errors.NewVerifyError()
	for _, cause := range vError.Causes {
		field := getErrorField(cause)
		msg := getErrorMsg(cause)
		err.AppendField(field, msg)
	}
	return err
}

func getErrorField(cause *jsonschema.ValidationError) string {
	field := strings.Join(cause.InstanceLocation, ".")
	return field
}

func getErrorMsg(cause *jsonschema.ValidationError) string {
	var kinds []string
	var errKind any = cause.ErrorKind
	switch errKind.(type) {
	case *kind.Format:
		kinds = append(kinds, "格式错误")
	case *kind.Required:
		var required = errKind.(*kind.Required)
		missing := strings.Join(required.Missing, ",")
		kinds = append(kinds, "必填项:"+missing)
	case *kind.Type:
		var kType = errKind.(*kind.Type)
		missing := strings.Join(kType.Want, ",")
		if len(missing) > 0 {
			missing = ",应为" + missing + "。"
		}
		kinds = append(kinds, "类型错误"+missing)
	case *kind.Enum:
		kinds = append(kinds, "枚举值错误")
	case *kind.MinLength:
		kinds = append(kinds, "长度过短")
	case *kind.MaxLength:
		kinds = append(kinds, "长度过长")
	case *kind.Minimum:
		kinds = append(kinds, "值过小")
	case *kind.Maximum:
		kinds = append(kinds, "值过大")
	case *kind.Pattern:
		kinds = append(kinds, "格式错误")
	case *kind.AdditionalItems:
		kinds = append(kinds, "额外项错误")
	case *kind.AdditionalProperties:
		kinds = append(kinds, "额外属性错误")
	case *kind.OneOf:
		kinds = append(kinds, "值不符合任何一个选项")
	case *kind.AnyOf:
		kinds = append(kinds, "值不符合任何一个选项")
	case *kind.AllOf:
		kinds = append(kinds, "值不符合所有选项")
	case *kind.Not:
		kinds = append(kinds, "值不应该出现")
	case *kind.InvalidJsonValue:
		kinds = append(kinds, "JSON格式错误")
	default:
		kinds = append(kinds, "未知错误")
	}
	return strings.Join(kinds, ",")
}
