package hugedao

import (
	"context"
	"encoding/json"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type Base[T any] struct {
	DBSchema *dbschema.DBSchema
	idao.Dao[T]
	client hugegraph.CommonClient
}

const DB_KEY_HUGE = "huge"

func (d *Base[T]) Write(ctx context.Context, script string, fb *stringutils.FmtBuilder, params map[string]any) error {
	hugeStore := d.GetStore()
	if hugeStore == nil {
		return nil
	}
	if fb != nil {
		script = fb.Format(script)
	}
	_, err := hugeStore.Write(ctx, script, params)
	if err != nil {
		logs.Error(ctx, logs.Fields{
			"err":    err,
			"script": script,
			"params": func() any {
				p, _ := json.Marshal(params)
				return string(p)
			},
		})
	} else {
		logs.InfoMsg(ctx, script)
	}
	return err
}

func (d *Base[T]) Query(ctx context.Context, gremlin string, fb *stringutils.FmtBuilder, params map[string]any) (*store_huge.HugeResult, error) {
	hugeStore := d.GetStore()
	if hugeStore == nil {
		return nil, errors.New("not found store")
	}
	return hugeStore.Query(ctx, gremlin, params)
}

func (d *Base[T]) GetStore() *store_huge.Dao[T] {
	iStore := d.Dao.GetStore().(any)
	hugeStore, ok := iStore.(*store_huge.Dao[T])
	if !ok {
		panic("hugegraph store does not implement store_huge.Dao")
	}
	return hugeStore
}

func (d *Base[T]) GetDBSchema() *dbschema.DBSchema {
	return d.DBSchema
}
