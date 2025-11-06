package ddd

type Entity interface {
	GetTenantId() string
	SetTenantId(v string)
	GetId() string
	SetId(v string)
}

type EntityList *[]Entity

/*
type MapEntity types.Object

func NewMapEntity() MapEntity {
	return MapEntity{}
}

func (e MapEntity) GetTenantId() string {
	return types.Object(e).GetString("tenantId")
}

func (e MapEntity) SetTenantId(val string) {
	_ = types.Object(e).Set("tenantId", val)
}

func (e MapEntity) GetId() string {
	return types.Object(e).GetString("id")
}

func (e MapEntity) SetId(val string) {
	_ = types.Object(e).Set("id", val)
}

func (e MapEntity) AsMap() map[string]any {
	return types.Object(e).AsMap()
}
*/
