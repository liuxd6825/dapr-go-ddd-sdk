package storage

import "context"

type RedisKeyValueStorage struct {
}

func NewRedisKeyValueStorage() *RedisKeyValueStorage {
	return &RedisKeyValueStorage{}
}

func (r *RedisKeyValueStorage) CreateTenant(ctx context.Context, tenantId string) error {
	return nil
}

func (r *RedisKeyValueStorage) CreateCase(ctx context.Context, tenantId, caseId string) error {
	return nil
}

func (r *RedisKeyValueStorage) DeleteTenant(ctx context.Context, tenantId string) error {
	return nil
}

func (r *RedisKeyValueStorage) DeleteCase(ctx context.Context, tenantId, caseId string) error {
	return nil
}
func (r *RedisKeyValueStorage) DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error {
	return nil
}

func (r *RedisKeyValueStorage) KVSource(ctx context.Context, id string) (Source, error) {
	return Source{}, nil
}

func (r *RedisKeyValueStorage) LoadTenant(ctx context.Context, tenantId, caseId string) error {
	return nil
}
func (r *RedisKeyValueStorage) KVUpsertSources(ctx context.Context, sources []Source) error {
	return nil
}
