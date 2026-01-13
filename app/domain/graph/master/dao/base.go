package dao

import (
	"context"
	"encoding/json"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type Base[T any] struct {
	DBSchema *dbschema.DBSchema
	idao.Dao[T]
}

func (d *Base[T]) write(ctx context.Context, fmtStr string, fb *stringutils.FmtBuilder, params map[string]any) error {
	nodeStore := d.GetStore()
	cypher := fb.Format(fmtStr)
	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		logs.Error(ctx, logs.Fields{"err": err, "cypher": cypher, "params": func() any {
			p, _ := json.Marshal(params)
			return string(p)
		}})

	} else {
		// logs.InfoMsg(ctx, cypher)
	}
	return err
}

func (d *Base[T]) GetStore() *store_neo4j.Dao[T] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[T])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
