package dao

import (
	"context"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

type IBusRelationDao interface {
	GraphMeta() *schema.Graph
	CreateSameName(ctx context.Context, node *model2.MasterNode)
	FindByName(ctx context.Context, name string) *model2.MasterRelation
	UpdateRecord(ctx context.Context, record *restapi.CDCRecord) *model2.MasterNode
	DeleteRecord(ctx context.Context, record *restapi.CDCRecord)
	GetStore() *store_huge.Dao[*model2.MasterRelation]
}
