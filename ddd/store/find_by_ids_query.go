package store

type FindByIdsQuery interface {
	GetIds() []string
	SetIds(val ...string) FindByIdsQuery

	GetTenantId() string
	SetTenantId(val string) FindByIdsQuery
}

type FindByIdsQueryRequest struct {
	TenantId string   `json:"tenantId" bson:"tenant_id"`
	Ids      []string `json:"ids" bson:"ids"`
}

func NewFindByIdsQuery() FindByIdsQuery {
	return &FindByIdsQueryRequest{}
}

func (q *FindByIdsQueryRequest) GetIds() []string {
	return q.Ids
}

func (q *FindByIdsQueryRequest) SetIds(val ...string) FindByIdsQuery {
	q.Ids = val
	return q
}

func (q *FindByIdsQueryRequest) GetTenantId() string {
	return q.TenantId
}

func (q *FindByIdsQueryRequest) SetTenantId(val string) FindByIdsQuery {
	q.TenantId = val
	return q
}
