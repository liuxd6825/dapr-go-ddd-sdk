package ddd

type Entity interface {
	GetTenantId() string
	GetId() string
}

var NilEntity = newNilEntity()

type EntityList *[]Entity

type nilEntity struct {
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
