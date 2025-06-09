package kv

type KVStorage interface {
	GetById(id string) (map[string]any, error)
	GetByIds(ids []string) ([]map[string]any, error)
}
