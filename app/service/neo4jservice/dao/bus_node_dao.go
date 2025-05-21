package dao

import (
	"context"
	"encoding/json"
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
	"strings"
)

type BusNodeDao struct {
	idao.Dao[*model.BusNode]
	*Base
	labels []string
}

type IStore[T any] interface {
	GetStore() store.IStore[T]
}

func NewBusNodeDao(labels []string, dbSch *store.DBSchema) *BusNodeDao {
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        labels,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	newDao := dao.NewDao[*model.BusNode](nodeCfg)
	return &BusNodeDao{
		Dao:    newDao,
		labels: labels,
		Base:   &Base{DBSchema: dbSch},
	}
}

func (d *BusNodeDao) Create(ctx context.Context, node *model.BusNode) {
	storeDao := d.GetStore()
	c := storeDao.Cypher
	props, dataMap, err := c.GetCreateProperties(ctx, node)
	if err != nil {
		panic(err)
	}
	tenantId := node.TenantId
	labels := c.GetLabels(ctx, node)
	cypher := fmt.Sprintf(`
	CREATE (n%s{%s}) WITH n 
	MERGE (m%s:same {name: $Name}) WITH n,m
	MERGE (n)-[:same]->(m)
	`, labels, props, labels)

	_, err = storeDao.DoSet(ctx, tenantId, cypher, dataMap)
	if err != nil {
		panic(err)
	}
}

// CreateRelNode 根据关系数据创建
func (d *BusNodeDao) CreateRelNode(ctx context.Context, rel *model.BusRelation, relNode *model.BusNode) *model.BusNode {
	storeDao := d.GetStore()
	c := storeDao.Cypher
	nLabels := storeDao.GetLabels(ctx, relNode)
	relType := rel.RelType

	nProps, data, err := c.GetCreateProperties(ctx, relNode)
	if err != nil {
		panic(err)
	}

	cypher := fmt.Sprintf(`
	CREATE (newNode%s{%s})  
	WITH newNode
	MATCH (mNode%s{id:$relStartId})
	CREATE (mNode)-[r:%s]->(newNode) 
	SET r.id=$relId, r.relType=$relType, r.case_id=$caseId, r.tenant_id=$tenantId;
`, nLabels, nProps, nLabels, relType)

	data["relId"] = rel.Id
	data["relType"] = rel.RelType
	data["relStartId"] = rel.StartId
	data["relEndId"] = rel.EndId
	data["tenantId"] = rel.TenantId
	data["caseId"] = rel.CaseId
	_, err = storeDao.Write(ctx, cypher, data)

	if err != nil {
		panic(err)
	}

	cypher2 := fmt.Sprintf(`
	MERGE (nName%s:same{name: $Name})      
	WITH nName         
	MATCH (n%s{id: $relId})          
	CREATE (n)-[r:same]->(nName) ;
	`, nLabels, nLabels)

	_, err = storeDao.Write(ctx, cypher2, data)

	if err != nil {
		panic(err)
	}
	return nil
}

func (d *BusNodeDao) Update(ctx context.Context, node *model.BusNode) {
	storeDao := d.GetStore()
	c := storeDao.Cypher

	props, dataMap, err := c.GetUpdateProperties(ctx, node, "n")
	if err != nil {
		panic(err)
	}

	labels := c.GetLabels(ctx, node)
	cypher := fmt.Sprintf(`
	MATCH (n%s) WHERE n.id=$Id SET %s WITH n 
	MERGE (n)-[:same]->(m%s:same{name:$Name}) WITH n
	MATCH (n)-[r:same]->(b) WHERE b.name <> $Name DELETE r
`, labels, props, labels)

	_, err = storeDao.Write(ctx, cypher, dataMap)
	if err != nil {
		panic(err)
	}
}

func (d *BusNodeDao) CreateByName(ctx context.Context, node *model.BusNode, name string) {
	storeDao := d.GetStore()
	c := storeDao.Cypher
	props, dataMap, err := c.GetUpdateProperties(ctx, node, "n")
	if err != nil {
		panic(err)
	}

	nodeLabels := c.GetLabels(ctx, node)

	//id := storeDao.GetId(node)
	cypher := fmt.Sprintf("MERGE (n%s{name:'%v'}) ON CREATE SET %v ON MATCH SET %v RETURN count(n) as rows ", nodeLabels, name, props, props)
	logs.Debug(ctx, logs.Fields{"cypher": cypher})
	_, err = storeDao.Write(ctx, cypher, dataMap)
	if err != nil {
		panic(err)
	}
}

func (d *BusNodeDao) FindByName(ctx context.Context, name string) *model.BusNode {
	nodes := d.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", name))
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// UpdateRelNode 根据关系数据创建图节点与关系
func (d *BusNodeDao) UpdateRelNode(ctx context.Context, record *model.Record) {
	// 是否改名
	isRename := record.IsRename()
	isChangedRelType := record.IsChangedRelType()
	if !isRename && !isChangedRelType {
		return
	}
	before := record.BeforeMap()
	after := record.AfterMap()
	node := model.NewBusNode(d.DBSchema.TableName, after)
	rel := model.NewBusRelation(d.DBSchema, after)
	oldName, _ := maputils.GetString(before, "name", "")

	nodeStore := d.GetStore()
	tenantId, _ := appctx.GetTenantId(ctx)
	relNodeLabel := nodeStore.GetLabels(ctx, node)

	params := map[string]interface{}{
		"tenantId": tenantId,
		"id":       node.Id,
		"name":     node.Name,
		"caseId":   node.CaseId,

		"relId":        rel.Id,
		"relType":      rel.RelType,
		"relTable":     rel.TableName,
		"startId":      rel.StartId,
		"oldNodeName":  oldName,
		"relNodeId":    node.Id,
		"relNodeName":  node.Name,
		"relNodeTable": node.Table,
	}

	sb := strings.Builder{}
	if isChangedRelType {
		sb.WriteString(fmt.Sprintf(`
			MATCH (n1%s)-[r1{id:$id}]->(m1%s) DELETE r1  WITH n1,m1
			CREATE (n1)-[r2:%s{id:$id}]->(m1) 
			`, relNodeLabel, relNodeLabel, rel.RelType,
		))
	}
	if isRename {
		sb.WriteString(fmt.Sprintf(`
		MATCH (n%s{id:$id}) SET n.name=$name WITH n
		MERGE (m%s:same{name:$name}) WITH n, m
		MERGE (n)-[r:same]->(m) WITH n
		MATCH (n)-[rd:same]->(md) WHERE md.name<>$name DELETE rd
		`, relNodeLabel, relNodeLabel))
	}

	_, err := nodeStore.Write(ctx, sb.String(), params)
	if err != nil {
		logs.Error(ctx, logs.Fields{"err": err, "cypher": sb.String(), "params": func() any {
			p, _ := json.Marshal(params)
			return string(p)
		}})
		panic(err)
	} else {
		logs.InfoMsg(ctx, sb.String())
	}

	return
}

func (d *BusNodeDao) DeleteRelNode(ctx context.Context, record *model.Record) {
	nodeStore := d.GetStore()
	after := record.AfterMap()

	rel := model.NewBusRelation(d.DBSchema, after)
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

	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *BusNodeDao) FindNodeAndRelationsByName(ctx context.Context, name string, labels ...string) *model.BusNode {
	storeDao := d.GetStore()
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

	//relCount := result.GetIntList("rel_count")
	for _, vMap := range n {
		node := d.NewNode(vMap)
		//node.RelCount = relCount[i]
		return node
	}

	return nil
}

func (d *BusNodeDao) ClearAll(ctx context.Context) {
	nodeStore := d.GetStore()
	cypher := fmt.Sprintf(`MATCH (n) DETACH DELETE n`)
	params := map[string]any{}
	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *BusNodeDao) NewNode(data map[string]any) *model.BusNode {
	node := &model.BusNode{
		Id:     d.GetString(data, "id"),
		CaseId: d.GetString(data, "case_id"),
		Name:   d.GetString(data, "name"),
		//RelCount: d.GetInt64(data, "rel_count"),
		TenantId: d.GetString(data, "tenant_id"),
		Table:    d.GetString(data, "table"),
	}
	return node
}

func (d *BusNodeDao) GetString(vMap map[string]any, dbName string) string {
	val, _ := maputils.GetString(vMap, dbName, "")
	return val
}

func (d *BusNodeDao) GetInt64(vMap map[string]any, dbName string) int64 {
	val, _ := maputils.GetInt64(vMap, dbName, -1)
	return val
}

func (d *BusNodeDao) getDriver() neo4j.DriverWithContext {
	driver := d.Dao.GetConfig().DB.(neo4j.DriverWithContext)
	return driver
}

func (d *BusNodeDao) GetStore() *store_neo4j.Dao[*model.BusNode] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model.BusNode])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
