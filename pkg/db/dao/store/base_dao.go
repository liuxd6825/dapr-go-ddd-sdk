package store

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
	return EntityToMap(entity, dao.schema, dao.eb.GetConfig().IsMap)
}

func (dao *BaseDao[T]) DB2Entity(data map[string]any) T {
	entity := dao.NewEntity()
	return MapToEntity[T](data, entity, dao.schema, dao.eb.GetConfig().IsMap)
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
	return ConvertValue(dataType, value)
}
