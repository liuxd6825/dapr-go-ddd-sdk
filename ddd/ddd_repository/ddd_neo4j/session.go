package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jSession struct {
	driver        neo4j.DriverWithContext
	sessionConfig neo4j.SessionConfig
	isWrite       bool
}

type SessionOptions struct {
	AccessMode *neo4j.AccessMode
}
type SessionOptionsBuilder struct {
	options SessionOptions
}

func NewSessionOptionsBuilder() *SessionOptionsBuilder {
	return &SessionOptionsBuilder{
		options: SessionOptions{},
	}
}

func (b *SessionOptionsBuilder) SetAccessMode(mode neo4j.AccessMode) *SessionOptionsBuilder {
	b.options.AccessMode = &mode
	return b
}

func (b *SessionOptionsBuilder) Build() *SessionOptions {
	return &SessionOptions{
		AccessMode: b.options.AccessMode,
	}
}

func NewSession(isWrite bool, driver neo4j.DriverWithContext) ddd_repository.Session {
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}
	if isWrite {
		sessionConfig.AccessMode = neo4j.AccessModeWrite
	}
	return &Neo4jSession{
		driver:        driver,
		sessionConfig: sessionConfig,
		isWrite:       isWrite,
	}
}

func (r *Neo4jSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc, opts ...*ddd_repository.SessionOptions) error {
	if dbFunc == nil {
		return nil
	}

	opt := ddd_repository.NewSessionOptions(opts...)
	session := r.driver.NewSession(ctx, r.sessionConfig)
	defer func() {
		_ = session.Close(ctx)
	}()

	txFunc := func(tx neo4j.ManagedTransaction) (res interface{}, err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		sessCtx := NewSessionContext(ctx, tx, session)
		err = dbFunc(sessCtx)
		return nil, err
	}

	var err error
	if r.isWrite {
		_, err = session.ExecuteWrite(ctx, txFunc, func(config *neo4j.TransactionConfig) {
			config.Timeout = opt.GetWriteTime()
		})
	} else {
		_, err = session.ExecuteRead(ctx, txFunc, func(config *neo4j.TransactionConfig) {
			config.Timeout = opt.GetReadTime()
		})
	}
	if err != nil {
		logs.Error(ctx, err)
	}
	return err
}

func (o *SessionOptions) SetAccessMode(accessMode neo4j.AccessMode) {
	o.AccessMode = &accessMode
}

func NewSessionOptions(opts ...*SessionOptions) *SessionOptions {
	s := &SessionOptions{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		if o.AccessMode != nil {
			s.SetAccessMode(*o.AccessMode)
		}
	}
	return s
}

func (o *SessionOptions) setDefault() {
	if o.AccessMode == nil {
		o.SetAccessMode(neo4j.AccessModeRead)
	}
}
