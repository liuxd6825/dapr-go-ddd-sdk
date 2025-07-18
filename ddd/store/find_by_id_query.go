package store

type FindByIdQuery interface {
	GetId() string
	SetId(val string) FindByIdQuery
}

type FindByIdQueryRequest struct {
	Id string `json:"id" path:"id" required:"true" ` // 聚合根Id
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
