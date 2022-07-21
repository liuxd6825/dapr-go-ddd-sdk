package ddd_neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type sessionKey struct {
}

type SessionContext interface {
	context.Context
	GetSession() neo4j.Session
	GetTransaction() neo4j.Transaction
}

func NewSessionContext(ctx context.Context, tr neo4j.Transaction, session neo4j.Session) SessionContext {
	return &sessionContext{
		Context: context.WithValue(ctx, sessionKey{}, session),
		tr:      tr,
		session: session,
	}
}

func GetSessionContext(ctx context.Context) (SessionContext, bool) {
	s := ctx.Value(sessionKey{})
	if s == nil {
		return nil, false
	}
	sess, ok := s.(SessionContext)
	return sess, ok
}

type sessionContext struct {
	context.Context
	tr      neo4j.Transaction
	session neo4j.Session
}

func (s *sessionContext) GetSession() neo4j.Session {
	return s.session
}

func (s *sessionContext) GetTransaction() neo4j.Transaction {
	return s.tr
}
