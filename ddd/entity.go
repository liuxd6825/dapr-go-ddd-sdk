package ddd

type Entity interface {
	GetTenantId() string
	GetId() string
}

type EntityList *[]Entity
