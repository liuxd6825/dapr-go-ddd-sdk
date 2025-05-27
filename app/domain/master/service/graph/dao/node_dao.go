package dao

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
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

type NodeDao struct {
	idao.Dao[*model2.Node]
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
	newDao := dao.NewDao[*model2.Node](nodeCfg)
	return &NodeDao{
		Dao:    newDao,
		labels: labels,
		Base:   &Base{DBSchema: dbSch},
	}
}

func (d *NodeDao) CreateMain(ctx context.Context, node *model2.Node) {
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
func (d *NodeDao) CreateRelNode(ctx context.Context, rel *model2.Relation, relNode *model2.Node) *model2.Node {
	storeDao := d.GetStore()
	c := storeDao.Cypher
	relType := rel.RelType
	nLabels := storeDao.GetLabels(ctx, relNode)
	mLabels := d.getLabels(rel.CaseId, rel.TenantId)
	nProps, data, err := c.GetCreateProperties(ctx, relNode)
	if err != nil {
		panic(err)
	}

	cypher := fmt.Sprintf(`
	CREATE (n%s{%s})   WITH n
	MATCH (m%s{id:$relStartId}) WITH m, n
	CREATE (m)-[r:%s]->(n) 
	SET r.id=$relId, r.relType=$relType, r.case_id=$caseId, r.tenant_id=$tenantId
    `, nLabels, nProps, mLabels, relType)

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

func (d *NodeDao) UpdateMain(ctx context.Context, node *model2.Node) {
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

// UpdateRelNode 根据关系数据创建图节点与关系
func (d *NodeDao) UpdateRelNode(ctx context.Context, record *model2.Record) {
	// 是否改名
	isRename := record.IsRename()
	isChangedRelType := record.IsChangedRelType()
	if !isRename && !isChangedRelType {
		return
	}
	before := record.BeforeMap()
	after := record.AfterMap()
	node := model2.NewNode(d.DBSchema.TableName, after)
	rel := model2.NewRelation(d.DBSchema, after)
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

func (d *NodeDao) DeleteMain(ctx context.Context, record *model2.Record) {
	nodeStore := d.GetStore()
	after := record.AfterMap()

	node := model2.NewNode(record.Table, after)
	labels := d.getLabels(node.CaseId, node.TenantId)

	cypher := fmt.Sprintf(`
	MATCH (n%s{id:$id})-[r]->(m%s) DETACH DELETE n,r,m  `, labels, labels)

	params := map[string]any{
		"id": node.Id,
	}

	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *NodeDao) DeleteRelNode(ctx context.Context, record *model2.Record) {
	nodeStore := d.GetStore()
	after := record.AfterMap()

	rel := model2.NewRelation(d.DBSchema, after)
	labels := d.getLabels(rel.CaseId, rel.TenantId)

	cypher := fmt.Sprintf(`
	MATCH (n%s)-[r{id:'%s'}]->(m%s) DETACH DELETE r,m`, labels, rel.Id, labels)

	_, err := nodeStore.Write(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
}

func (d *NodeDao) getLabels(caseId, tenantId string, label ...string) string {
	labels := fmt.Sprintf(":case_%s:tenant_%s:master", caseId, tenantId)
	for _, l := range label {
		labels += fmt.Sprintf(":%s", l)
	}
	return labels
}

func (d *NodeDao) ClearAll(ctx context.Context) {
	nodeStore := d.GetStore()
	cypher := fmt.Sprintf(`MATCH (n) DETACH DELETE n`)
	params := map[string]any{}
	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *NodeDao) NewNode(data map[string]any) *model2.Node {
	node := &model2.Node{
		Id:     d.GetString(data, "id"),
		CaseId: d.GetString(data, "case_id"),
		Name:   d.GetString(data, "name"),
		//RelCount: d.GetInt64(data, "rel_count"),
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

func (d *NodeDao) GetStore() *store_neo4j.Dao[*model2.Node] {
	iStore := d.Dao.GetStore().(any)
	nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model2.Node])
	if !ok {
		panic("neo4j store does not implement neo4j.Dao")
	}
	return nodeStoreDao
}
