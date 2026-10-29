package huge

import (
	"context"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type Table struct {
	cfg    store2.DBSchema
	client *hugegraph.CommonClient
}

func newTable(client *hugegraph.CommonClient, schema *store2.DBSchema) *Table {
	return &Table{client: client, cfg: *schema}
}

func (t *Table) GetTableName() string {
	return t.cfg.TableName
}

func (t *Table) GetSchema() *store2.DBSchema {
	return &t.cfg
}

func (t *Table) AutoMigrate(ctx context.Context) error {
	/*	return store_huge.SyncSchema(ctx, t.client, store_huge.SchemaSpec{

		})*/
	return nil

}

func (t *Table) Exist(ctx context.Context) bool {
	return t.client != nil
}

func (t *Table) Drop(ctx context.Context) error {
	_, err := t.client.Gremlin.Post(
		t.client.Gremlin.Post.WithGremlin("g.V().drop().iterate();g.E().drop().iterate()"),
	)
	return err
}
