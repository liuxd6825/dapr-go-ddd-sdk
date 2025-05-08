package db_pkg

import "context"

type Table interface {
	Name() string
	Create(ctx context.Context) (err error)
	Exist(ctx context.Context) (bool, error)
	Drop(ctx context.Context) error
}
