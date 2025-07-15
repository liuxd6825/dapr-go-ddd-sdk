package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	ParamTag    = "param"    // 路径/查询参数标签
	BodyTag     = "body"     // 请求体标签
	PathTag     = "path"     // 路径参数标签
	QueryTag    = "query"    // 查询参数标签
	RequiredTag = "required" // 必填参数标签
	JsonTag     = "json"     // json数据体
	TitleTag    = "title"    // 字段标题标签
	ValidateTag = "validate"
)

// GetParams 将请求参数绑定到目标结构体
func GetParams(ctx iris.Context, target interface{}) (err error) {
	// 处理请求体JSON
	if ctx.Request().Method == iris.MethodPost || ctx.Request().Method == iris.MethodPut {
		err = ctx.ReadJSON(target)
		if err != nil {
			return err
		}
		err = Validate(target)
		return err
	}

	// 数据验证对象
	verifyErr := errors.NewVerifyError()

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetType := targetValue.Elem().Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Elem().Field(i)

		if pathName := field.Tag.Get(PathTag); pathName != "" {
			// 处理路径参数 /api/v1.0/user/{id}
			if val := ctx.Params().Get(pathName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if queryName := field.Tag.Get(QueryTag); queryName != "" {
			// 处理查询参数 ?name=lxd
			if val := ctx.URLParam(queryName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if paramName := field.Tag.Get(ParamTag); paramName != "" {
			// 处理通用param标签（兼容双模式）
			if val := ctx.Params().Get(paramName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			} else if val := ctx.URLParam(paramName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if field.Type.Kind() == reflect.Struct { // 处理嵌套结构体
			err := bindNestedStruct(ctx, &fieldValue)
			if err != nil {
				// 检查是否为VerifyError类型，合并子结构体的验证错误
				if ve, ok := err.(*errors.VerifyError); ok {
					verifyErr.MergeWithPrefix(field.Name, ve)
				} else {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		}
	}

	return verifyErr.GetError()
}

func Validate(target interface{}) error {

	verifyError := errors.NewVerifyError()
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetType := targetValue.Elem().Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Elem().Field(i)
		validField(verifyError, &field, fieldValue)
	}
	return verifyError.GetError()
}

// validField 递归校验单个字段，支持结构体、结构体指针、必填校验
func validField(verifyError *errors.VerifyError, field *reflect.StructField, fieldValue reflect.Value) {
	// 支持 time.Time 及其别名类型（如 times.Time）必填校验
	if isTimeType(field.Type) {
		isRequired := fieldIsRequired(field)
		if isRequired && isZeroTime(fieldValue) {
			appendFieldError(verifyError, field, "不能为空")
		}
		return
	}
	// 支持 *time.Time 及其别名类型指针必填校验
	if isTimePtrType(field.Type) {
		isRequired := fieldIsRequired(field)
		if isRequired && (fieldValue.IsNil() || isZeroTime(fieldValue.Elem())) {
			appendFieldError(verifyError, field, "不能为空")
		}
		return
	}
	// 结构体类型递归
	if field.Type.Kind() == reflect.Struct {
		err := Validate(fieldValue.Addr().Interface())
		if ve, ok := err.(*errors.VerifyError); ok && ve.Count() > 0 {
			verifyError.MergeWithPrefix(getJsonFieldName(field), ve)
		}
		return
	}
	// 结构体指针类型递归和必填
	if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
		required := field.Tag.Get(RequiredTag)
		if strings.ToLower(required) == "true" && fieldValue.IsNil() {
			appendFieldError(verifyError, field, "不能为空")
			return
		}
		if !fieldValue.IsNil() {
			err := Validate(fieldValue.Interface())
			if ve, ok := err.(*errors.VerifyError); ok && ve.Count() > 0 {
				verifyError.MergeWithPrefix(getJsonFieldName(field), ve)
			}
		}
		return
	}

	// 普通字段必填
	required := field.Tag.Get(RequiredTag)
	if strings.ToLower(required) == "true" && (fieldValue.IsZero() || (fieldValue.Kind() == reflect.String && fieldValue.String() == "")) {
		appendFieldError(verifyError, field, "不能为空")
	}
}

func fieldIsRequired(field *reflect.StructField) bool {
	required := field.Tag.Get(RequiredTag)
	if strings.ToLower(required) == "true" {
		return true
	}
	validate := field.Tag.Get(ValidateTag)
	if strings.Contains(validate, "required") {
		return true
	}
	return false
}

func appendFieldError(verifyError *errors.VerifyError, field *reflect.StructField, msg string) {
	name := getJsonFieldName(field)
	title := field.Tag.Get(TitleTag)
	verifyError.AppendField(name, msg, title)
}

func getJsonFieldName(field *reflect.StructField) string {
	name := field.Tag.Get(JsonTag)
	if name == "" {
		name = stringutils.FirstLower(field.Name)
	}
	return name
}

// bindNestedStruct 处理嵌套结构体绑定
func bindNestedStruct(ctx iris.Context, field *reflect.Value) error {
	nestedPtr := reflect.New(field.Type())
	if err := GetParams(ctx, nestedPtr.Interface()); err != nil {
		return err
	}
	field.Set(nestedPtr.Elem())
	return nil
}

// setFieldValue 设置字段值并处理类型转换
func setFieldValue(field *reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	}
	return nil
}

// 判断是否为 time.Time 或 times.Time 及其别名类型
func isTimeType(t reflect.Type) bool {
	if t == reflect.TypeOf(time.Time{}) || t == reflect.TypeOf(times.Time{}) {
		return true
	}
	if t == reflect.TypeOf(times.Time{}) || t == reflect.TypeOf(times.Time{}) {
		return true
	}
	if t == reflect.TypeOf(times.Date{}) || t == reflect.TypeOf(times.Date{}) {
		return true
	}
	return false
}

// 判断是否为 *time.Time 或 *times.Time 及其别名类型
func isTimePtrType(t reflect.Type) bool {
	if t == reflect.TypeOf(&time.Time{}) || t == reflect.TypeOf(&times.Time{}) {
		return true
	}
	if t == reflect.TypeOf(&times.Time{}) || t == reflect.TypeOf(&times.Time{}) {
		return true
	}
	if t == reflect.TypeOf(&times.Date{}) || t == reflect.TypeOf(&times.Date{}) {
		return true
	}
	return false
}

// 判断时间是否为零值
func isZeroTime(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	if v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}
	if v.Type() == reflect.TypeOf(times.Time{}) {
		return v.Interface().(times.Time) == times.Time{}
	}
	if v.Type() == reflect.TypeOf(times.Date{}) {
		return v.Interface().(times.Date) == times.Date{}
	}
	return true
}
