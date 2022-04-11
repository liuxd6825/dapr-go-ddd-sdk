package ddd_repository

import "context"

type Session interface {
	UseTransaction(context.Context, SessionFunc) error
}

type SessionFunc func(ctx context.Context) error

func StartSession(ctx context.Context, session Session, dbFunc SessionFunc) error {
	return session.UseTransaction(ctx, dbFunc)
}
