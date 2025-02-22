package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
)

type Table struct {
	tableName string
	schema    *jsonschema.Schema
	db        *ddd_mongodb.MongoDB
}

func NewTable(db *ddd_mongodb.MongoDB, schema *jsonschema.Schema) db.Table {
	return newTable(db, schema)
}

func newTable(db *ddd_mongodb.MongoDB, schema *jsonschema.Schema) *Table {
	tableName := stringutils.AsFieldName(schema.Name)
	return &Table{db: db, tableName: tableName, schema: schema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *jsonschema.Schema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	var err error
	if ctx == nil {
		ctx = context.Background()
	}
	if isExist, e := t.db.ExistCollection(ctx, t.tableName); e != nil {
		err = e
	} else if !isExist {
		err = t.db.CreateCollection(t.tableName)
	}
	if err != nil {
		panic(err)
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	isExist, err := t.db.ExistCollection(ctx, t.tableName)
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
	if isExist, e := t.db.ExistCollection(ctx, t.tableName); e != nil {
		err = e
	} else if isExist {
		err = t.db.GetCollection(t.tableName).Drop(ctx)
	}
	if err != nil {
		panic(err)
	}
}
