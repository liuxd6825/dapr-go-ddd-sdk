package store

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/convert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapToEntity[T any](data map[string]any, entityAny T, schema *DBSchema, entityIsMap bool) T {
	if entityIsMap {
		var anyData any
		anyData = entityAny
		eMap, isMap := anyData.(map[string]any)
		if !isMap {
			panic("entity is not a map")
		}
		for _, field := range schema.Fields {
			val := data[field.DBName]
			if field.DataType == DataType_Date || field.DataType == DataType_Time {
				if pDate, ok := val.(primitive.DateTime); ok {
					timeVal := pDate.Time().In(times.GetLocalTimeZone())
					val = &timeVal
				}
			} else if field.DataType == DataType_Array {
				if arr, ok := val.(primitive.A); ok {
					timeVal := arr
					val = &timeVal
				}
			}
			eMap[field.Name] = val
		}
	} else {
		var errFieldName string
		gp.Try(func() error {
			for _, field := range schema.Fields {
				errFieldName = field.Name
				val := data[field.DBName]
				if field.DataType == DataType_Date || field.DataType == DataType_Time {
					if pDate, ok := val.(primitive.DateTime); ok {
						timeVal := pDate.Time()
						val = &timeVal
					}
				} else if field.DataType == DataType_Array {
					if arr, ok := val.(primitive.A); ok {
						timeVal := arr
						val = &timeVal
					}
				}
				if val != nil {
					err := reflectutils.SetField(entityAny, field.Name, val)
					if err != nil {
						return errors.New("set field %s to entity %T error: %s", field.Name, entityAny, err.Error())
					}
				}
			}
			return nil
		}).Catch(func(err error) {
			panic(fmt.Sprintf("转换db值到属性%s时出错：%s", errFieldName, err.Error()))
		})
	}
	return entityAny
}

// EntityToMap 将实体对象转化为map[string]any
func EntityToMap(entity any, schema *DBSchema, entityIsMap bool) map[string]any {
	data := map[string]any{}
	eMap, isMap := entity.(map[string]any)
	if entityIsMap {
		if !isMap {
			panic("entity is not map")
		}
	}

	if isMap {
		if entityIsMap {
			for _, field := range schema.Fields {
				fieldValue, ok := eMap[field.Name]
				if !ok && field.Name != field.DBName {
					fieldValue, ok = eMap[field.DBName]
				}
				if !ok {
					//fieldValue = field.DefaultValueInterface
				}

				val, err := ConvertValue(field.DataType, fieldValue)
				if err != nil {
					panic(fmt.Sprintf("fieldName:%s, dataType:%s, value: %v, error:%s", field.Name, DataTypeToString(field.DataType), fieldValue, err.Error()))
				}
				data[field.DBName] = val
			}
		} else {
			newMap := map[string]any{}
			for key, fieldValue := range eMap {
				field := schema.LookedField(key)
				if field != nil {
					val, err := ConvertValue(field.DataType, fieldValue)
					if err != nil {
						panic(fmt.Sprintf("fieldName:%s, dataType:%s, value: %v, error:%s", field.Name, DataTypeToString(field.DataType), fieldValue, err.Error()))
					}
					newMap[field.DBName] = val
				}
			}
			return newMap
		}
	} else {
		for _, field := range schema.Fields {
			val := reflectutils.GetField(entity, field.Name)
			data[field.DBName] = val
		}
	}
	return data
}

func ConvertValue(dataType DataType, value any) (any, error) {
	if value == nil {
		return value, nil
	}
	switch dataType {
	case DataType_Bool:
		return convert.ToBool(value)
	case DataType_String:
		convert.ToString(value)
	case DataType_Int:
		return convert.ToInt(value)
	case DataType_Float:
		return convert.ToFloat64(value)
	case DataType_Time:
		return convert.ToTime(value)
	case DataType_Date:
		return convert.ToTime(value)
	}
	return value, nil
}
