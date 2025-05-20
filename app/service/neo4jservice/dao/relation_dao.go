package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type RelationDao struct {
	idao.Dao[*model.Relation]
	*Base
}

func NewRelationDao(dbSch *store.DBSchema) *RelationDao {
	relCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Rel,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	relDao := dao.NewDao[*model.Relation](relCfg)
	return &RelationDao{
		Dao:  relDao,
		Base: &Base{DBSchema: dbSch},
	}
}

func (d *RelationDao) FindByName(ctx context.Context, name string) *model.Relation {
	nodes := d.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", name))
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}
