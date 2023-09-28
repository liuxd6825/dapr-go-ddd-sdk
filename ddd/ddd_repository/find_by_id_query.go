package ddd_repository

type IFindByIdQuery interface {
	GetId() string
	SetId(val string) FindByIdQuery

	GetTenantId() string
	SetTenantId(val string) FindByIdQuery
}

type FindByIdQuery struct {
	TenantId string `json:"tenantId" bson:"tenant_id"`
	Id       string `json:"id" bson:"id"`
}

func NewFindByIdQuery() *FindByIdQuery {
	return &FindByIdQuery{}
}

func (q *FindByIdQuery) GetId() string {
	return q.Id
}

func (q *FindByIdQuery) SetId(val string) *FindByIdQuery {
	q.Id = val
	return q
}

func (q *FindByIdQuery) GetTenantId() string {
	return q.TenantId
}

func (q *FindByIdQuery) SetTenantId(val string) *FindByIdQuery {
	q.TenantId = val
	return q
}
