package db

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
)

type Db struct {
}

func NewDb() *Db {
	return &Db{}
}
func Open() (*Db, error) {
	return &Db{}, nil
}

func (db *Db) Model(schema *schema.Schema) *Model {
	return NewModel(schema, db)
}

func (db *Db) First(ctx context.Context, schema *schema.Schema, val ...string) Record {
	return NewModel(schema, db).First(ctx, val...)
}

func (db *Db) Create(ctx context.Context, schema *schema.Schema, record Record) error {
	return NewModel(schema, db).Create(ctx, record)
}

func (db *Db) Delete(ctx context.Context, schema *schema.Schema, val ...string) error {
	return NewModel(schema, db).Delete(ctx, val...)
}
