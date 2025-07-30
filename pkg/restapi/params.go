package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/validator"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
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

// GetWebParams 将请求参数绑定到目标结构体
func GetWebParams(ictx iris.Context, target interface{}, removeNames ...string) (res any, err error) {
	// 处理请求体JSON
	request := ictx.Request()
	if request.Method == iris.MethodPost || request.Method == iris.MethodPut {
		if request.ContentLength == 0 {
			return nil, errors.New("request body is null")
		}
		err = ictx.ReadJSON(target)
		if err != nil {
			return nil, err
		}
		if dataMap, ok := target.(*map[string]any); ok {
			return dataMap, err
		}
		if dataMap, ok := target.(map[string]any); ok {
			return dataMap, err
		}
		err = validator.Validate(target, removeNames...)
		return res, err
	}

	// 数据验证对象
	verifyErr := errors.NewVerifyError()

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("target must be a pointer to a struct")
	}

	targetType := targetValue.Elem().Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Elem().Field(i)
		if pathName := field.Tag.Get(PathTag); pathName != "" {
			// 处理路径参数 /api/v1.0/user/{id}
			if val := ictx.Params().Get(pathName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if queryName := field.Tag.Get(QueryTag); queryName != "" {
			// 处理查询参数 ?name=lxd
			if val := ictx.URLParam(queryName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if paramName := field.Tag.Get(ParamTag); paramName != "" {
			// 处理通用param标签（兼容双模式）
			if val := ictx.Params().Get(paramName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			} else if val := ictx.URLParam(paramName); val != "" {
				if err := setFieldValue(&fieldValue, val); err != nil {
					verifyErr.AppendField(field.Name, err.Error())
				}
			}
		} else if field.Type.Kind() == reflect.Struct { // 处理嵌套结构体
			err := bindNestedStruct(ictx, &fieldValue)
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
	if err := verifyErr.IsHasError(); err {
		return nil, verifyErr.GetError()
	}
	res = targetValue.Interface()
	return res, validator.Validate(target, removeNames...)
}

// bindNestedStruct 处理嵌套结构体绑定
func bindNestedStruct(ctx iris.Context, field *reflect.Value) error {
	nestedPtr := reflect.New(field.Type())
	param, err := GetWebParams(ctx, nestedPtr.Interface())
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(param))
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
