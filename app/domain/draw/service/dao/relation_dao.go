package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type RelationDao struct {
	idao.Dao[*model.Relation]
}

func NewRelationDao() *RelationDao {
	relCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Rel,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	relDao := dao.NewDao[*model.Relation](relCfg)
	return &RelationDao{
		Dao: relDao,
	}
}

func (d *RelationDao) GetStore() *store_neo4j.Dao[*model.Relation] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model.Relation])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
