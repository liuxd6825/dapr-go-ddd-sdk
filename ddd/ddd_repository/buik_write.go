package ddd_repository

type BulkWriteResult struct {
	InsertedCount int64
	MatchedCount  int64
	ModifiedCount int64
	DeletedCount  int64
	UpsertedCount int64
	UpsertedIDs   map[int64]any
	isEmpty       bool
}

func NewBulkWriteResult() *BulkWriteResult {
	return &BulkWriteResult{isEmpty: true}
}

func (r *BulkWriteResult) IsEmpty() bool {
	if r == nil {
		return true
	}
	return r.isEmpty
}

func (r *BulkWriteResult) SetEmpty(val bool) {
	r.isEmpty = val
}
