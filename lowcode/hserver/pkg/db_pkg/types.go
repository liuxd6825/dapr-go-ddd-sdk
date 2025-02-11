package db_pkg

// AccessType DAO数据访问类型
type AccessType string

const (
	AccessTypeCreate AccessType = "create"
	AccessTypeUpdate AccessType = "update"
	AccessTypeDelete AccessType = "delete"
)
