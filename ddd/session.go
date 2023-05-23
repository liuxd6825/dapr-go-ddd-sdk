package ddd

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type Session interface {
	GetSessionId() string
	GetTenantId() string
}

type session struct {
	tenantId  string
	sessionId string
}

type sessionKey struct {
}

type eventItem struct {
	aggregate     Aggregate
	event         DomainEvent
	ctx           context.Context
	callEventType CallEventType
	opts          []*ApplyEventOptions
}

func StartSession(ctx context.Context, tenantId string, do func(ctx context.Context, session Session) error) (resErr error) {
	defer func() {
		if err := errors.GetRecoverError(recover()); err != nil {
			resErr = err
		}
	}()
	session, ok := getSession(ctx)
	if !ok {
		ctx, session = newContext(ctx, tenantId)
	}
	err := do(ctx, session)
	if err != nil {
		if e := session.rollback(ctx); e != nil {
			return errors.ErrorOf("事件回滚失败,详细:%s。回滚原因:%s", e.Error(), err.Error())
		}
		return err
	}
	return session.commit(ctx)
}

func newContext(parent context.Context, tenantId string) (context.Context, *session) {
	ctx, _ := context.WithCancel(parent)
	value := newSession(tenantId)
	return context.WithValue(ctx, sessionKey{}, value), value
}

func newSession(tenantId string) *session {
	value := &session{
		tenantId:  tenantId,
		sessionId: uuid.New().String(),
	}
	return value
}

func getSession(ctx context.Context) (*session, bool) {
	value := ctx.Value(sessionKey{})
	events, ok := value.(*session)
	return events, ok
}

func (b *session) GetTenantId() string {
	return b.tenantId
}

func (b *session) GetSessionId() string {
	return b.sessionId
}

func (b *session) commit(ctx context.Context) error {
	_, err := Commit(ctx, b.tenantId, b.sessionId)
	return err
}

func (b *session) rollback(ctx context.Context) error {
	_, err := Rollback(ctx, b.tenantId, b.sessionId)
	return err
}
