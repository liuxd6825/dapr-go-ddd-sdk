package ddd_repository

type FindByIdQuery interface {
	TenantId() string
	Id() string
}

type findByIdQuery struct {
	tenantId string
	id       string
}

func NewFindByIdQuery(tenantId, id string) FindByIdQuery {
	return &findByIdQuery{
		tenantId: tenantId,
		id:       id,
	}
}

func (q *findByIdQuery) TenantId() string {
	return q.tenantId
}

func (q *findByIdQuery) Id() string {
	return q.id
}
