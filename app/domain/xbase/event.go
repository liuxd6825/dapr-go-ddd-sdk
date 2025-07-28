package xbase

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
	"strings"
	"sync"
)

var _outboxDao idao.OutboxEventDao
var _outboxOnce sync.Once

func PublishEvent(ctx context.Context, appId string, data any, meta map[string]any) error {
	tenantId, err := appctx.GetTenantId3(ctx)
	if err != nil {
		return err
	}
	outboxDao := newOutboxEventDao()
	event, err := getEvent(ctx, appId, tenantId, data, meta)
	if err != nil {
		return err
	}
	return outboxDao.Create(ctx, event).GetError()
}

func newOutboxEventDao() idao.OutboxEventDao {
	_outboxOnce.Do(func() {
		_outboxDao = dao.NewOutboxEventDao(config.DBKey)
	})
	return _outboxDao
}

func getEvent(ctx context.Context, appId string, tenantId string, data any, meta map[string]any) (*dbevent.OutboxEvent, error) {
	var err error
	_, topic, _ := reflectutils.GetTypeDetails(data)
	topic = stringutils.MidlineString(topic)

	dataMap, err := convertStructToMapViaReflection(data)
	if err != nil {
		return nil, err
	}
	event := &dbevent.OutboxEvent{
		Id:          idutils.NewId(),
		AppId:       appId,
		TenantId:    tenantId,
		Topic:       topic,
		Data:        dataMap,
		Meta:        meta,
		CreatedTime: times.NewTime(),
	}
	return event, err
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
