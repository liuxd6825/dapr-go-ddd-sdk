package store

type FindByIdQuery interface {
	GetId() string
	SetId(val string) FindByIdQuery

	GetTenantId() string
	SetTenantId(val string) FindByIdQuery
}

type FindByIdQueryRequest struct {
	TenantId string `json:"tenantId" `
	Id       string `json:"id" path:"id" required:"true" ` // 聚合根Id
}

func NewFindByIdQuery() FindByIdQuery {
	return &FindByIdQueryRequest{}
}

func (q *FindByIdQueryRequest) GetId() string {
	return q.Id
}

func (q *FindByIdQueryRequest) SetId(val string) FindByIdQuery {
	q.Id = val
	return q
}

func (q *FindByIdQueryRequest) GetTenantId() string {
	return q.TenantId
}

func (q *FindByIdQueryRequest) SetTenantId(val string) FindByIdQuery {
	q.TenantId = val
	return q
}
