package db_pkg

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type Table struct {
	name   string
	schema *schema.Schema
	db     *DB
}

func NewTable(db *DB, name string, schema *schema.Schema) *Table {
	return &Table{db: db, name: name, schema: schema}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Create(ctx context.Context) (err error) {
	if isExist, e := t.db.mongodb.ExistCollection(ctx, t.name); e != nil {
		err = e
	} else if !isExist {
		err = t.db.mongodb.CreateCollection(t.name)
	}
	return err
}

func (t *Table) Exist(ctx context.Context) (bool, error) {
	isExist, err := t.db.mongodb.ExistCollection(ctx, t.name)
	return isExist, err
}

func (t *Table) Drop(ctx context.Context) error {
	if isExist, err := t.db.mongodb.ExistCollection(ctx, t.name); err != nil {
		return err
	} else if isExist {
		return t.db.mongodb.GetCollection(t.name).Drop(ctx)
	}
	return nil
}
