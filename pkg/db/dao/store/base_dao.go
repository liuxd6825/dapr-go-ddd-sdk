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

type BaseDao[T any] struct {
	eb     EntityBuilder[T] // 实体构造器
	schema *DBSchema
}

func NewBaseDao[T any](eb EntityBuilder[T], schema *DBSchema) *BaseDao[T] {
	return &BaseDao[T]{
		eb:     eb,
		schema: schema,
	}
}

// Entity2DB 将实体对象转化为map[string]any
func (dao *BaseDao[T]) Entity2DB(entity any) map[string]any {
	data := map[string]any{}
	entityIsMap := dao.eb.GetConfig().IsMap
	eMap, isMap := entity.(map[string]any)
	if entityIsMap {
		if !isMap {
			panic("entity is not map")
		}
	}

	if isMap {
		if entityIsMap {
			for _, field := range dao.schema.Fields {
				fieldValue, ok := eMap[field.Name]
				if !ok && field.Name != field.DBName {
					fieldValue, ok = eMap[field.DBName]
				}
				if !ok {
					//fieldValue = field.DefaultValueInterface
				}

				val, err := dao.ConvertValue(field.DataType, fieldValue)
				if err != nil {
					panic(fmt.Sprintf("Dao.ConvertValue() fileName:%s, dataType:%s, error:%s", field.Name, DataTypeToString(field.DataType), err.Error()))
				}
				data[field.DBName] = val
			}
		} else {
			newMap := map[string]any{}
			for key, fieldValue := range eMap {
				field := dao.schema.LookedField(key)
				if field != nil {
					val, err := dao.ConvertValue(field.DataType, fieldValue)
					if err != nil {
						panic(fmt.Sprintf("Dao.ConvertValue() fileName:%s, dataType:%s, error:%s", field.Name, DataTypeToString(field.DataType), err.Error()))
					}
					newMap[field.DBName] = val
				}
			}
			return newMap
		}
	} else {
		for _, field := range dao.schema.Fields {
			val := reflectutils.GetField(entity, field.Name)
			data[field.DBName] = val
		}
	}
	return data
}

func (dao *BaseDao[T]) DB2Entity(data map[string]any) T {
	entity := dao.NewEntity()
	isMap := dao.eb.GetConfig().IsMap
	if isMap {
		eAny := any(entity)
		eMap, isMap := eAny.(map[string]any)
		if !isMap {
			panic("store_mongodb.dao entity is not a map")
		}
		for _, field := range dao.GetSchema().Fields {
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
			for _, field := range dao.GetSchema().Fields {
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
					err := reflectutils.SetField(entity, field.Name, val)
					if err != nil {
						return errors.New("set field %s to entity %T error: %s", field.Name, entity, err.Error())
					}
				}
			}
			return nil
		}).Catch(func(err error) {
			panic(fmt.Sprintf("转换db值到属性%s时出错：%s", errFieldName, err.Error()))
		})
	}
	return entity
}

func (dao *BaseDao[T]) GetSchema() *DBSchema {
	return dao.schema
}

func (dao *BaseDao[T]) NewEntity() T {
	return dao.eb.NewEntity()
}

func (dao *BaseDao[T]) NewEntityList() []T {
	return dao.eb.NewEntityList()
}

func (dao *BaseDao[T]) GetTenantId(entity T) string {
	return dao.eb.GetTenantId(entity)
}

func (dao *BaseDao[T]) SetTenantId(entity T, tenantId string) {
	dao.eb.SetTenantId(entity, tenantId)
}

func (dao *BaseDao[T]) GetId(entity T) string {
	return dao.eb.GetId(entity)
}

func (dao *BaseDao[T]) SetId(entity T, id string) {
	dao.eb.SetId(entity, id)
}

func (dao *BaseDao[T]) GetAggId(entity T) string {
	return dao.eb.GetAggId(entity)
}

func (dao *BaseDao[T]) GetEb() EntityBuilder[T] {
	return dao.eb
}

func (dao *BaseDao[T]) ConvertValue(dataType DataType, value any) (any, error) {
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
