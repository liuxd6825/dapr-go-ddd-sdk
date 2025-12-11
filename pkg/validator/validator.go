package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate, _ = New()
}

func New() (*validator.Validate, error) {
	// 初始化验证器
	validate = validator.New()
	// 注册 TagNameFunc：优先使用 label 标签，否则字段名
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		if label := f.Tag.Get("title"); label != "" {
			return label
		}
		return f.Name
	})

	// 初始化中文翻译器
	zhCn := zh.New()
	uni := ut.New(zhCn, zhCn)
	trans, _ = uni.GetTranslator("zh")

	// 注册中文翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func Validate(target interface{}) error {
	err := validate.Struct(target)
	return GetVerifyError(err)
}

// ValidateExcept
// @Description: 验证传入的字段之外的所有字段。
// @param target
// @param fields
// @return error
func ValidateExcept(target interface{}, fields ...string) error {
	err := validate.StructExcept(target, fields...)
	return GetVerifyError(err)
}

// ValidatePartial
// @Description:  验证传入的字段
// @param target
// @param fields
// @return error
func ValidatePartial(target interface{}, fields ...string) error {
	err := validate.StructPartial(target, fields...)
	return GetVerifyError(err)
}

func GetVerifyError(err error) error {
	if err != nil {
		verifyError := errors.NewVerifyError()
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			if fieldErr, ok := err.(validator.FieldError); ok {
				title := getTitle(fieldErr)
				structNamespace := fieldErr.StructNamespace()
				sn := strings.Split(structNamespace, ".")
				count := len(sn)
				if count > 0 {
					title = sn[count-1]
					sn = sn[1:]
				}
				msg := fieldErr.Translate(trans)
				field := getField(fieldErr, nil)
				verifyError.AppendField(field, msg, title)
			}
		}
		return verifyError.GetError()
	}
	return nil
}

func getTitle(fieldErr validator.FieldError) string {
	names := strings.Split(fieldErr.Namespace(), ".")
	return stringutils.FirstLower(names[len(names)-1])
}

func getField(fieldErr validator.FieldError, cancelFields []string) string {
	names := strings.Split(fieldErr.StructNamespace(), ".")
	var res []string
	for i, name := range names {
		if i == 0 {
			continue
		}

		if stringutils.Include(name, cancelFields, true) {
			continue
		}
		name = stringutils.FirstLower(name)
		res = append(res, name)
	}
	return strings.Join(res, ".")
}
