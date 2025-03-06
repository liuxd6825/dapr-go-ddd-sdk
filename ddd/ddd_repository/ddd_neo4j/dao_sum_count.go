package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

func (d *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) map[string]any {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}
