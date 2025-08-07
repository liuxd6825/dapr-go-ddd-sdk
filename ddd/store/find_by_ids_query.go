package store

type FindByIdsQuery interface {
	GetIds() []string
	SetIds(val ...string) FindByIdsQuery
}

type FindByIdsQueryRequest struct {
	Ids []string `json:"ids" query:"ids" validate:"-"`
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
