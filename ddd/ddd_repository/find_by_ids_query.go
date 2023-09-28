package ddd_repository

type IFindByIdsQuery interface {
	GetIds() []string
	SetIds(val ...string) IFindByIdsQuery

	GetTenantId() string
	SetTenantId(val string) IFindByIdsQuery
}

type FindByIdsQuery struct {
	TenantId string   `json:"tenantId" bson:"tenant_id"`
	Ids      []string `json:"ids" bson:"ids"`
}

func NewFindByIdsQuery() *FindByIdsQuery {
	return &FindByIdsQuery{}
}

func (q *FindByIdsQuery) GetIds() []string {
	return q.Ids
}

func (q *FindByIdsQuery) SetIds(val ...string) *FindByIdsQuery {
	q.Ids = val
	return q
}

func (q *FindByIdsQuery) GetTenantId() string {
	return q.TenantId
}

func (q *FindByIdsQuery) SetTenantId(val string) *FindByIdsQuery {
	q.TenantId = val
	return q
}
