package ddd

type Entity interface {
	GetTenantId() string
	SetTenantId(v string)
	GetId() string
	SetId(v string)
}

var NilEntity = newNilEntity()

type EntityList *[]Entity

type nilEntity struct {
}

func (e *nilEntity) SetTenantId(v string) {
	//TODO implement me
	panic("implement me")
}

func (e *nilEntity) SetId(v string) {
	//TODO implement me
	panic("implement me")
}

func newNilEntity() Entity {
	return &nilEntity{}
}

func (e *nilEntity) GetTenantId() string {
	return ""
}

func (e *nilEntity) GetId() string {
	return ""
}
