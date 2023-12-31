package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jSession struct {
	driver        neo4j.Driver
	sessionConfig neo4j.SessionConfig
	isWrite       bool
}

func NewSession(isWrite bool, driver neo4j.Driver) ddd_repository.Session {
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

func (r *Neo4jSession) UseTransaction(ctx context.Context, dbFunc ddd_repository.SessionFunc) error {
	session := r.driver.NewSession(r.sessionConfig)
	defer func() {
		_ = session.Close()
	}()

	tx := func(tx neo4j.Transaction) (res interface{}, resErr error) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					resErr = err
				}
			}
		}()
		sessCtx := NewSessionContext(ctx, tx, session)
		if dbFunc == nil {
			return nil, nil
		}
		err := dbFunc(sessCtx)
		if err != nil {
			return nil, tx.Rollback()
		}
		return nil, tx.Commit()
	}

	var err error
	if r.isWrite {
		_, err = session.WriteTransaction(tx)
	} else {
		_, err = session.ReadTransaction(tx)
	}
	return err
}
