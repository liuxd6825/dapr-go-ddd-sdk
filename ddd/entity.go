package ddd

type Entity interface {
	GetTenantId() string
	GetId() string
}

var NilEnity = newNilEnity()

type EntityList *[]Entity

type nilEnity struct {
}

func (e *nilEnity) GetTenantId() string {
	return ""
}
func (e *nilEnity) GetId() string {
	return ""
}

func newNilEnity() Entity {
	return &nilEnity{}
}
