package vector

import "context"

type InsertData struct {
	TenantId string
	CaseId   string
	Source   string
	DocId    string
	Batch    []string
	Vectors  [][]float32
}

type VectorStorage interface {
	//Init(ctx context.Context, collectionName string, dim int) error
	Insert(ctx context.Context, Insert InsertData) error
	DropCollection(ctx context.Context) error
	Flush(ctx context.Context) error
	Search(ctx context.Context, vector []float32, topK int) ([]string, error)
	/*
		CreateVectorIndexIfNotExists()
		Upsert(data map[string]map[string]any)
		Query(query string, topK int64, ids []string) []map[string]any
		IndexDoneCallback()
		Delete(ids []string) error
		DeleteEntity(entityName string) error
		DeleteEntityRelation(entityName string) error
		GetById(id string) (map[string]any, error)
		GetByIds(ids []string) ([]map[string]any, error)
		Drop() map[string]any
	*/
}
