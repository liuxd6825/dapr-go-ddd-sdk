package events

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"reflect"
	"strings"
)

// GetEventParams 接受事件端取得事件参数
func GetEventParams(cloudEvent CloudEvent, target interface{}) (res any, ctx context.Context, err error) {
	cloudData, err := cloudEvent.GetData()
	if err != nil {
		return nil, nil, err
	}

	var payload EventPayload
	err = json.Unmarshal([]byte(cloudData), &payload)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.Background()
	if payload.Data == nil {
		return target, ctx, nil
	}
	err = json.Unmarshal(payload.Data, target)
	if err != nil {
		return nil, nil, err
	}

	if payload.Meta != nil {
		var meta map[string]any
		err = json.Unmarshal(payload.Meta, &meta)
		if err != nil {
			return nil, nil, err
		}
		if val, ok := meta["autoUser"]; ok {
			meta[AUTH_USER] = val
		}

		if mapVal, ok := meta[AUTH_USER]; ok {
			if authUserMap, ok := mapVal.(map[string]any); ok {
				authUser := appctx.NewAuthUserEntity(authUserMap)
				ctx, err = appctx.NewAuthUserContext(ctx, authUser)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		if idVal, ok := meta[TENANT_ID]; ok {
			if tenantId, ok := idVal.(string); ok {
				ctx = appctx.NewTenantContext(ctx, tenantId)
			}
		}

		if headerVal, ok := meta[HEADER]; ok {
			if header, ok := headerVal.(map[string][]string); ok {
				ctx = appctx.NewHeaderContext(ctx, header)
			}
		}
	}
	return target, ctx, err
}

func newMeta(ctx context.Context, tenantId string, meta map[string]any) (map[string]any, error) {
	if meta == nil {
		meta = make(map[string]any)
	}
	if _, ok := meta[TENANT_ID]; !ok {
		meta[TENANT_ID] = tenantId
	}

	if _, ok := meta[AUTH_USER]; !ok {
		if autoUser, ok := appctx.GetAuthUser(ctx); ok {
			autoUserMap, err := convertStructToMapViaReflection(autoUser)
			if err != nil {
				return nil, err
			}
			meta[AUTH_USER] = autoUserMap
		}
	}

	if _, ok := meta[HEADER]; !ok {
		if header, ok := appctx.GetHeader(ctx); ok {
			meta[HEADER] = header
		} else {
			meta[HEADER] = appctx.Header{}
		}
	}
	return meta, nil
}

// getFieldNameFromTag 从字段的 tag 中获取用作 map key 的名称，优先使用 json tag
func getFieldNameFromTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return field.Name
	}
	// `json:"name,omitempty"` -> "name"
	return strings.Split(tag, ",")[0]
}

// ConvertStructToMapViaReflection 使用反射将结构体转为Map
func convertStructToMapViaReflection(data any) (map[string]any, error) {
	result := make(map[string]any)
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// 如果是指针，需要获取其指向的元素
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 确保是结构体类型
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("input data must be a struct or a pointer to a struct")
	}

	// 辅助函数，用于递归处理
	var processFields func(v reflect.Value, t reflect.Type)
	processFields = func(v reflect.Value, t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// 如果是匿名嵌入字段，则递归处理
			if field.Anonymous {
				// 确保嵌入的字段也是一个结构体
				if fieldValue.Kind() == reflect.Struct {
					processFields(fieldValue, fieldValue.Type())
				}
				continue
			}

			// 对于普通字段，获取其 tag 作为 key，并存入 map
			fieldName := getFieldNameFromTag(field)
			result[fieldName] = fieldValue.Interface()
		}
	}

	processFields(val, typ)
	return result, nil
}
