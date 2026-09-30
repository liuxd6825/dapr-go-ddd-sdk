package dao

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

type IMasterNodeDao interface {
	idao.Dao[*model.MasterNode]
	GraphMeta() *schema.Graph
	CreateMain(ctx context.Context, node *model.MasterNode)
	CreateRelNode(ctx context.Context, rel *model.MasterRelation, relNode *model.MasterNode) *model.MasterNode
	UpdateMain(ctx context.Context, node *model.MasterNode)
	UpdateRelNode(ctx context.Context, record *dbevent.CDCRecord) error
	DeleteMain(ctx context.Context, record *dbevent.CDCRecord, dbSch *dbschema.DBSchema)
	DeleteRelNode(ctx context.Context, record *dbevent.CDCRecord)
	ClearAll(ctx context.Context)
	NewNode(data map[string]any) *model.MasterNode
	GetString(vMap map[string]any, dbName string) string
	GetInt64(vMap map[string]any, dbName string) int64

	IsRename(r *dbevent.CDCRecord) bool
	IsChangedRelType(r *dbevent.CDCRecord) bool
	GetSchema() *dbschema.DBSchema
	GetConfig() *idao.DaoConfig
}

type IStore[T any] interface {
	GetStore() store2.IStore[T]
}
