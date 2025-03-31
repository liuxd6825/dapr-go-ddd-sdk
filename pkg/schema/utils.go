package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/jsonschema/v6"
)

func GetTableName(sch *jsonschema.Schema) string {
	tableName := sch.Name()
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			if meta.DBTable != nil && meta.DBTable.Name != "" {
				tableName = meta.DBTable.Name
				break
			}
		}
	}
	return tableName
}

func GetFieldName(sch *jsonschema.Schema) string {
	fieldName := sch.Name()
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			if meta.DBField != nil && meta.DBField.Name != "" {
				fieldName = meta.DBField.Name
				break
			}
		}
	}
	return fieldName
}

func GetField(sch *jsonschema.Schema) *DBField {
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			return meta.DBField
		}
	}
	return nil
}

func GetMetaExtension(sch *jsonschema.Schema) *MetaExtension {
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			return meta
		}
	}
	return nil
}

func ApplyDefaults(schema *jsonschema.Schema, data any) interface{} {
	switch v := data.(type) {
	case map[string]any:
		return ApplyDefaultsToMap(schema, v)
	case []interface{}:
		return ApplyDefaultsToArray(schema, v)
	default:
		return data
	}
}

func ApplyDefaultsToMap(sch *jsonschema.Schema, data map[string]any) any {
	if sch.Properties != nil {
		props := sch.GetAllProperties()
		for propName, propSchema := range props {
			if _, exists := data[propName]; !exists {
				if propSchema.Default != nil {
					data[propName] = propSchema.Default
				} else {
					// 根据类型设置空值
					if propSchema.Types.Contains(jsonschema.JsonType_ObjectType) {
						data[propName] = ApplyDefaultsToMap(propSchema, make(map[string]interface{}))
					} else if propSchema.Types.Contains(jsonschema.JsonType_StringType) {
						data[propName] = ApplyDefaultsToArray(propSchema, []interface{}{})
					}
				}
			} else {
				// 递归处理已有值
				data[propName] = ApplyDefaults(propSchema, data[propName])
			}
		}
	}
	return data
}

func ApplyDefaultsToStruct(sch *jsonschema.Schema, data any) any {
	if sch.Properties != nil {
		props := sch.GetAllProperties()
		refObj := reflectutils.NewRefObj(data)
		for propName, propSchema := range props {
			field := refObj.Field(propName)
			if field.IsValid() {
				val, err := field.Get()
				if err != nil {
					panic(err)
				}
				if propSchema.Default != nil && val == nil {
					err = field.Set(propSchema.Default)
				}
				if err != nil {
					panic(err)
				}
			} else {
				val, err := field.Get()
				if err != nil {
					panic(err)
				}
				// 递归处理已有值
				err = field.Set(ApplyDefaults(propSchema, val))
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return data
}

func ApplyDefaultsToArray(sch *jsonschema.Schema, data []interface{}) []interface{} {
	if sch.Items != nil {
		for i, item := range data {
			if s, ok := sch.Items.(*jsonschema.Schema); ok {
				data[i] = ApplyDefaults(s, item)
			} else if items, ok := item.([]*jsonschema.Schema); ok {
				for i, item := range items {
					data[i] = ApplyDefaults(item, data[i])
				}
			}
		}
	}
	return data
}
