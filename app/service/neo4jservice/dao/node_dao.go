package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type NodeDao struct {
	idao.Dao[*model.Node]
	*Base
	labels []string
}

type IStore[T any] interface {
	GetStore() store.IStore[T]
}

func NewNodeDao(labels []string, dbSch *store.DBSchema) *NodeDao {
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        labels,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	newDao := dao.NewDao[*model.Node](nodeCfg)
	return &NodeDao{
		Dao:    newDao,
		labels: labels,
		Base:   &Base{DBSchema: dbSch},
	}
}

func (d *NodeDao) MergeByName(ctx context.Context, node *model.Node, name string, labels ...string) *model.Node {
	storeDao := d.getNeo4jDao()
	c := storeDao.Cypher
	props, dataMap, err := c.GetUpdateProperties(ctx, node, "n")
	if err != nil {
		panic(err)
	}

	nodeLabels := c.GetLabels(ctx, node)

	id := storeDao.GetId(node)
	cypher := fmt.Sprintf("MERGE (n%s{id:'%v'}) ON CREATE SET %v ON MATCH SET %v RETURN count(n) as rows ", nodeLabels, id, props, props)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	_, err = storeDao.Write(ctx, cypher, dataMap)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d *NodeDao) FindByName(ctx context.Context, name string) *model.Node {
	nodes := d.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", name))
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// CreateByRelationData 根据关系数据创建
func (d *NodeDao) CreateByRelationData(ctx context.Context, rel *model.Relation, relNode *model.Node) *model.Node {
	storeDao := d.getNeo4jDao()
	nLabels := storeDao.GetLabels(ctx, relNode)
	mLabels := storeDao.GetLabels(ctx, relNode)
	relType := rel.RelType

	cypher := fmt.Sprintf(`
	MERGE (n%s{id: $nId})
	MERGE (m%s{name: $mName})
		ON CREATE SET m.id=$relNodeId, m.name=$relNodeName, m.table=$relNodeTable, m.tenant_id=$tenantId
	MERGE (n)-[r:%s]->(m)
	    SET r.id = $relId, r.relType=$relType, r.case_id=$relCaseId, r.tenant_id=$tenantId`,
		nLabels, mLabels, relType)
	params := map[string]interface{}{
		"tenantId":     relNode.TenantId,
		"nId":          rel.StartId,
		"mName":        relNode.Name,
		"relId":        rel.Id,
		"relType":      rel.RelType,
		"relCaseId":    rel.CaseId,
		"relNodeId":    relNode.Id,
		"relNodeName":  relNode.Name,
		"relNodeTable": relNode.Table,
	}

	_, err := storeDao.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
	return nil
}

// MergeAndCreateRel 根据关系数据创建图节点与关系
func (d *NodeDao) MergeAndCreateRel(ctx context.Context, node *model.Node, rel *model.Relation, relNode *model.Node) *model.Node {

	storeDao := d.getNeo4jDao()
	tenantId, _ := appctx.GetTenantId(ctx)
	startNodeLabel := storeDao.GetLabels(ctx, node)
	relNodeLabel := storeDao.GetLabels(ctx, relNode)
	cypher := fmt.Sprintf(`
	MERGE (n%s{id: $startId})   
		ON CREATE SET n.id=$id, n.name=$name, n.case_id=$caseId, n.tenant_id=$tenantId
	WITH n
		MERGE (m%s {name: $relNodeName})   
		ON CREATE SET m.tenant_id=$tenantId,m.id=$relNodeId,m.name=$relNodeName,m.table=$relTable

	MERGE (n)-[r%s{id:$relId}]->(m)
		ON CREATE SET r.id=$relId,r.tenant_id=$tenantId,r.case_id=$caseId,r.rel_type=$relType
	RETURN n, m;`, startNodeLabel, relNodeLabel, rel.RelType)
	params := map[string]interface{}{
		"tenantId": tenantId,
		"id":       node.Id,
		"name":     node.Name,
		"caseId":   node.CaseId,

		"relId":    rel.Id,
		"relType":  rel.RelType,
		"relTable": rel.TableName,
		"startId":  rel.StartId,

		"relNodeId":    relNode.Id,
		"relNodeName":  relNode.Name,
		"relNodeTable": relNode.Table,
	}
	_, err := storeDao.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
	return nil
}

func (d *NodeDao) DeleteNodeAndRelationByRecord(ctx context.Context, record *model.Record) {
	storeDao := d.getNeo4jDao()
	after := record.AfterMap()

	rel := model.NewRelation(d.DBSchema, after)
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

func (d *NodeDao) FindNodeAndRelationsByName(ctx context.Context, name string, labels ...string) *model.Node {
	storeDao := d.getNeo4jDao()
	nodeLabel := storeDao.GetLabels(ctx, nil, labels...)

	cypher := fmt.Sprintf("MATCH (n%s{name:'%s'}) OPTIONAL MATCH (n)-[r]->() RETURN n, count(r) AS rel_count", nodeLabel, name)
	result, err := storeDao.Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	dbSch := d.Dao.GetSchema()
	var n []map[string]any
	err = result.GetList(ctx, "n", &n, dbSch)
	if err != nil {
		panic(err)
	}

	relCount := result.GetIntList("rel_count")
	for i, vMap := range n {
		node := d.NewNode(vMap)
		node.RelCount = relCount[i]
		return node
	}

	return nil
}

func (d *NodeDao) NewNode(data map[string]any) *model.Node {
	node := &model.Node{
		Id:       d.GetString(data, "id"),
		CaseId:   d.GetString(data, "case_id"),
		Name:     d.GetString(data, "name"),
		RelCount: d.GetInt64(data, "rel_count"),
		TenantId: d.GetString(data, "tenant_id"),
		Table:    d.GetString(data, "table"),
	}
	return node
}

func (d *NodeDao) GetString(vMap map[string]any, dbName string) string {
	val, _ := maputils.GetString(vMap, dbName, "")
	return val
}

func (d *NodeDao) GetInt64(vMap map[string]any, dbName string) int64 {
	val, _ := maputils.GetInt64(vMap, dbName, -1)
	return val
}

func (d *NodeDao) getDriver() neo4j.DriverWithContext {
	driver := d.Dao.GetConfig().DB.(neo4j.DriverWithContext)
	return driver
}

func (d *NodeDao) getNeo4jDao() *store_neo4j.Dao[*model.Node] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model.Node])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
