package dao

import (
	"context"
	"encoding/json"
	"fmt"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
)

type BusRelationDao struct {
	idao.Dao[*model2.MasterRelation]
	nodeDao *MasterNodeDao
	*Base[*model2.MasterRelation]
	graphMeta *schema.Graph
}

func NewBusRelationDao(dbSch *store.DBSchema, nodeDao *MasterNodeDao, graphMeta *schema.Graph) *BusRelationDao {
	relCfg := &dao.DaoConfig{
		DBKey:              "neo4j",
		GraphType:          idao.GraphType_Rel,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	relDao := dao.NewDao[*model2.MasterRelation](relCfg)
	return &BusRelationDao{
		graphMeta: graphMeta,
		Dao:       relDao,
		nodeDao:   nodeDao,
		Base:      &Base[*model2.MasterRelation]{DBSchema: dbSch},
	}
}

func (d *BusRelationDao) GraphMeta() *schema.Graph {
	return d.graphMeta
}

func (d *BusRelationDao) CreateSameName(ctx context.Context, node *model2.MasterNode) {

}

func (d *BusRelationDao) FindByName(ctx context.Context, name string) *model2.MasterRelation {
	nodes, err := d.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", name))
	if err != nil {
		panic(err)
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// UpdateRecord 根据关系数据创建图节点与关系
func (d *BusRelationDao) UpdateRecord(ctx context.Context, record *restapi.CDCRecord) *model2.MasterNode {
	after := record.AfterMap()
	node := model2.NewMasterNode(after, record.DBSchema)
	rel := model2.NewMasterRelation(after, record.DBSchema)
	nodeName, _ := maputils.GetString(record.BeforeMap(), "name", "")
	nodeStore := d.nodeDao.GetStore()
	relStore := d.GetStore()
	tenantId, _ := appctx.GetTenantId(ctx)
	startNodeLabel := nodeStore.GetLabels(ctx, node)
	relNodeLabel := nodeStore.GetLabels(ctx, node)

	cypher := fmt.Sprintf(`
	MERGE (n%s{id: $source})   
		ON CREATE SET n.id=$id, n.name=$name
	WITH n
		MERGE (m%s {name: $relNodeName})   
		ON CREATE SET m.tenant_id=$tenantId,m.name=$relNodeName, m.description=$description
	MERGE (n)-[r:%s{id:$relId}]->(m)
		ON CREATE SET r.id=$relId,r.tenant_id=$tenantId,r.case_id=$caseId,r.rel_type=$relType, m.description=$description
	RETURN n, m;`, startNodeLabel, relNodeLabel, rel.RelType)
	params := map[string]interface{}{
		"tenantId":     tenantId,
		"id":           node.Id,
		"name":         node.Name,
		"caseId":       node.CaseId,
		"relId":        rel.Id,
		"relType":      rel.RelType,
		"relTable":     rel.SourceIds,
		"source":       rel.Source,
		"oldNodeName":  nodeName,
		"relNodeId":    node.Id,
		"relNodeName":  node.Name,
		"relNodeTable": node.SourceIds,
		"description":  rel.Description,
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

func (d *BusRelationDao) DeleteRecord(ctx context.Context, record *restapi.CDCRecord) {
	storeDao := d.GetStore()
	after := record.AfterMap()

	rel := model2.NewMasterRelation(after, d.DBSchema)
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

func (d *BusRelationDao) GetStore() *store_neo4j.Dao[*model2.MasterRelation] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model2.MasterRelation])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
