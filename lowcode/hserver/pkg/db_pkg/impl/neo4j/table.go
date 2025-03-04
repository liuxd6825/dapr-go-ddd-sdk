package neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Table struct {
	tableName string
	schema    *jsonschema.Schema
	driver    neo4j.DriverWithContext
}

func NewTable(driver neo4j.DriverWithContext, schema *jsonschema.Schema) db.Table {
	return newTable(driver, schema)
}

func newTable(driver neo4j.DriverWithContext, schema *jsonschema.Schema) *Table {
	tableName := stringutils.AsFieldName(schema.Name)
	return &Table{driver: driver, tableName: tableName, schema: schema}
}

func (t *Table) GetTableName() string {
	return t.tableName
}

func (t *Table) GetSchema() *jsonschema.Schema {
	return t.schema
}

func (t *Table) AutoMigrate(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
}

func (t *Table) Exist(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return true
}

func (t *Table) Drop(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

}
