package ddd_repository

import "context"

type SessionFunc func(ctx context.Context) error

type DbSession interface {
	UseSession(ctx context.Context, fn SessionFunc) error
}
