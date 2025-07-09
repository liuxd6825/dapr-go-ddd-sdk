package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	paramTag = "param" // 路径/查询参数标签
	bodyTag  = "body"  // 请求体标签
	pathTag  = "path"  // 路径参数标签
	queryTag = "query" // 查询参数标签
)

// GetParams 将请求参数绑定到目标结构体
func GetParams(ctx iris.Context, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	targetType := targetValue.Elem().Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Elem().Field(i)

		// 处理路径参数
		if pathName := field.Tag.Get(pathTag); pathName != "" {
			if val := ctx.Params().Get(pathName); val != "" {
				if err := setFieldValue(fieldValue, val); err != nil {
					return err
				}
				continue
			}
		}

		// 处理查询参数
		if queryName := field.Tag.Get(queryTag); queryName != "" {
			if val := ctx.URLParam(queryName); val != "" {
				if err := setFieldValue(fieldValue, val); err != nil {
					return err
				}
				continue
			}
		}

		// 处理通用param标签（兼容旧版）
		if paramName := field.Tag.Get(paramTag); paramName != "" {
			if val := ctx.Params().Get(paramName); val != "" {
				if err := setFieldValue(fieldValue, val); err != nil {
					return err
				}
				continue
			}
			if val := ctx.URLParam(paramName); val != "" {
				if err := setFieldValue(fieldValue, val); err != nil {
					return err
				}
				continue
			}
		}

		// 处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			if err := bindNestedStruct(ctx, fieldValue); err != nil {
				return err
			}
		}
	}

	// 处理请求体JSON
	/*
		if err := ctx.ReadJSON(target); err != nil && !iris.IsErrPath(err) {
			return err
		}
	*/
	return nil
}

// bindNestedStruct 处理嵌套结构体绑定
func bindNestedStruct(ctx iris.Context, field reflect.Value) error {
	nestedPtr := reflect.New(field.Type())
	if err := GetParams(ctx, nestedPtr.Interface()); err != nil {
		return err
	}
	field.Set(nestedPtr.Elem())
	return nil
}

// setFieldValue 设置字段值并处理类型转换
func setFieldValue(field reflect.Value, value string) error {
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
