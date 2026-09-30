package hugedao

import (
	"context"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	store_huge "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

// EnsureSchema 触发 schema 同步(add-only);
// SyncSchema 自身已幂等,无需本地 cache。
func EnsureSchema(ctx context.Context, client *hugegraph.CommonClient, graphMeta *schema.Graph) error {
	if client == nil || graphMeta == nil {
		return nil
	}
	return store_huge.SyncMasterSchema(ctx, client)
}
