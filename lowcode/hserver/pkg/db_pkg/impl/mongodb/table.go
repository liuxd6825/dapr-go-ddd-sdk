package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type Table struct {
	name   string
	schema *schema.Schema
	db     *ddd_mongodb.MongoDB
}

func NewTable(db *ddd_mongodb.MongoDB, name string, schema *schema.Schema) *Table {
	return &Table{db: db, name: name, schema: schema}
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) GetSchema() *schema.Schema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	var err error
	if ctx == nil {
		ctx = context.Background()
	}
	if isExist, e := t.db.ExistCollection(ctx, t.name); e != nil {
		err = e
	} else if !isExist {
		err = t.db.CreateCollection(t.name)
	}
	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	isExist, err := t.db.ExistCollection(ctx, t.name)
	if err != nil {
		panic(err)
	}
	return isExist
}

func (t *Table) Drop(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	if isExist, e := t.db.ExistCollection(ctx, t.name); e != nil {
		err = e
	} else if isExist {
		err = t.db.GetCollection(t.name).Drop(ctx)
	}
	if err != nil {
		panic(err)
	}
}
