package dao

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type BusRelationDao struct {
	idao.Dao[*model2.Relation]
	nodeDao *NodeDao
	*Base
}

func NewBusRelationDao(dbSch *store.DBSchema, nodeDao *NodeDao) *BusRelationDao {
	relCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Rel,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	relDao := dao.NewDao[*model2.Relation](relCfg)
	return &BusRelationDao{
		Dao:     relDao,
		nodeDao: nodeDao,
		Base:    &Base{DBSchema: dbSch},
	}
}

func (d *BusRelationDao) CreateSameName(ctx context.Context, node *model2.Node) {

}

func (d *BusRelationDao) FindByName(ctx context.Context, name string) *model2.Relation {
	nodes := d.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", name))
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// UpdateRecord 根据关系数据创建图节点与关系
func (d *BusRelationDao) UpdateRecord(ctx context.Context, record *model2.Record) *model2.Node {
	after := record.AfterMap()
	node := model2.NewNode(d.DBSchema.TableName, after)
	rel := model2.NewRelation(d.DBSchema, after)
	nodeName, _ := maputils.GetString(record.BeforeMap(), "name", "")
	nodeStore := d.nodeDao.GetStore()
	relStore := d.GetStore()
	tenantId, _ := appctx.GetTenantId(ctx)
	startNodeLabel := nodeStore.GetLabels(ctx, node)
	relNodeLabel := nodeStore.GetLabels(ctx, node)

	cypher := fmt.Sprintf(`
	MERGE (n%s{id: $startId})   
		ON CREATE SET n.id=$id, n.name=$name
	WITH n
		MERGE (m%s {name: $relNodeName})   
		ON CREATE SET m.tenant_id=$tenantId,m.name=$relNodeName
	MERGE (n)-[r:%s{id:$relId}]->(m)
		ON CREATE SET r.id=$relId,r.tenant_id=$tenantId,r.case_id=$caseId,r.rel_type=$relType
	RETURN n, m;`, startNodeLabel, relNodeLabel, rel.RelType)
	params := map[string]interface{}{
		"tenantId": tenantId,
		"id":       node.Id,
		"name":     node.Name,
		"caseId":   node.CaseId,

		"relId":        rel.Id,
		"relType":      rel.RelType,
		"relTable":     rel.TableName,
		"startId":      rel.StartId,
		"oldNodeName":  nodeName,
		"relNodeId":    node.Id,
		"relNodeName":  node.Name,
		"relNodeTable": node.Table,
	}
	_, err := relStore.Write(ctx, cypher, params)
	if err != nil {
		logs.Error(ctx, logs.Fields{"err": err, "cypher": cypher, "params": func() any {
			p, _ := json.Marshal(params)
			return string(p)
		}})
		panic(err)
	}
	return nil
}

func (d *BusRelationDao) DeleteRecord(ctx context.Context, record *model2.Record) {
	storeDao := d.GetStore()
	after := record.AfterMap()

	rel := model2.NewRelation(d.DBSchema, after)
	name, _ := maputils.GetString(after, "name", "")
	labels := fmt.Sprintf(":case_%s:tentant_%s:master", rel.CaseId, rel.TenantId)

	cypher := fmt.Sprintf(`MATCH (n%s)-[r:%s]->(m%s)
	WHERE r.id = $id AND (n.name = $name OR m.name = $name)
	WITH n, m, r, 
		 CASE 
			 WHEN n.name = $name THEN n
			 WHEN m.name = $name THEN m
		 END AS nodeToDelete
	MATCH (nodeToDelete)-[rel]->()  
	WITH nodeToDelete, r, COUNT(rel) AS relCount
	WHERE relCount = 1  
	DELETE nodeToDelete, r`, labels, rel.RelType, labels)

	params := map[string]any{
		"name": name,
		"id":   rel.Id,
	}

	_, err := storeDao.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *BusRelationDao) GetStore() *store_neo4j.Dao[*model2.Relation] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model2.Relation])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
