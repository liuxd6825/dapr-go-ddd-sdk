package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type GraphDao struct {
	dao   idao.Dao[map[string]any]
	store *store_neo4j.Dao[map[string]any]
}

func NewGraphDao() *GraphDao {
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"master"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
		DBSchema:           dbschema.NewDBSchema("graph", "graph"),
	}
	newDao := dao.NewDao[map[string]any](nodeCfg)
	return &GraphDao{
		dao: newDao,
	}
}

func (d *GraphDao) FindByCaseId(ctx context.Context, caseId string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:master:case_%s) OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) FindById(ctx context.Context, caseId string, id string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:master:case_%s{id:%s}) OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId, id)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) FindByName(ctx context.Context, caseId string, name string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:master:case_%s) WHERE n.name='%s' OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId, name)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) FindByContainName(ctx context.Context, caseId string, name string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:master:case_%s) WHERE n.name CONTAINS '%s' OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId, name)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) FindByStartWithName(ctx context.Context, caseId string, name string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:master:case_%s) WHERE n.name STARTS WITH '%s' OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId, name)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) GetStore() *store_neo4j.Dao[map[string]any] {
	if d.store == nil {
		iStore := d.dao.GetStore().(any)
		nodeStoreDao, ok := iStore.(*store_neo4j.Dao[map[string]any])
		if !ok {
			panic("neo4j store does not implement neo4j.Dao")
		}
		d.store = nodeStoreDao
	}
	return d.store
}
