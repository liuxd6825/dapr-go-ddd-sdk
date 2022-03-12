package ddd

type Entity interface {
	GetTenantId() string
	GetId() string
}

type View interface {
	GetTenantId() string
	GetId() string
}
