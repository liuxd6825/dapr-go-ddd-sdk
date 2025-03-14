package idao

// AccessType DAO数据访问类型
type AccessType string

const (
	AccessTypeCreate AccessType = "create"
	AccessTypeUpdate AccessType = "update"
	AccessTypeDelete AccessType = "delete"

	AccessTypeBatchCreate AccessType = "batch-create"
	AccessTypeBatchUpdate AccessType = "batch-update"
	AccessTypeBatchDelete AccessType = "batch-delete"
)
