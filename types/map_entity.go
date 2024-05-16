package types

type MapEntity Object

const (
	TenantId = "tenantId"
	Id       = "id"
)

func NewEntity() MapEntity {
	return MapEntity{}
}

func (e MapEntity) GetTenantId() string {
	return Object(e).GetString(TenantId)
}

func (e MapEntity) SetTenantId(val string) {
	_ = Object(e).Set(TenantId, val)
}

func (e MapEntity) GetId() string {
	return Object(e).GetString(Id)
}

func (e MapEntity) SetId(val string) {
	_ = Object(e).Set(Id, val)
}
