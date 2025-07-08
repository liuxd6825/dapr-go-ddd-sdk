package dao

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	labels := c.GetLabels(ctx, node)
	fmtStr := `
	CREATE (n$<labels>{$<props>}) WITH n 
	MERGE (m$<labels>:same {id:$<name>,name:$<name>}) WITH n,m
	MERGE (n)-[:same]->(m)`

	fb := stringutils.NewFmtBuilder()
	fb.String("labels", labels)
	fb.String("props", props)
	fb.Varchar("name", node.Name)

	err = d.write(ctx, fmtStr, fb, dataMap)
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

	fmtStr := `
	CREATE (n$<nLabels>{$<props>}) WITH n
	MATCH (m$<mLabels>{id:$<source>}) WITH m, n
	CREATE (m)-[r:$<relType>]->(n) 
	SET r.id=$<id>,r.keywords=$<keywords>,r.case_id=$<caseId>,r.tenant_id=$<tenantId>, 
		r.description=$<description>,r.source=$<source>,r.target=$<target>, 
		r.source_ids=$<id>,r.source_type=$<sourceType>,r.table=$<table>
    `

	fb := stringutils.NewFmtBuilder()
	fb.String("nLabels", nLabels)
	fb.String("props", nProps)
	fb.String("mLabels", mLabels)
	fb.String("relType", relType)
	fb.Varchar("id", rel.Id)
	fb.Varchar("source", rel.Source)
	fb.Varchar("target", rel.Target)
	fb.Varchar("tenantId", rel.TenantId)
	fb.Varchar("caseId", rel.CaseId)
	fb.Varchar("keywords", rel.Keywords)
	fb.Varchar("description", "")
	fb.Varchar("sourceType", rel.SourceType)
	fb.Varchar("table", rel.Table)

	_, err = storeDao.Write(ctx, fb.Format(fmtStr), data)

	if err != nil {
		panic(err)
	}

	fmtStr2 := `
	MERGE (nName$<labels>:same{name: $<name>}) WITH nName         
	MATCH (n$<labels>{id: $<id>})          
	CREATE (n)-[r:same]->(nName);
	`

	fb2 := stringutils.NewFmtBuilder()
	fb2.String("labels", nLabels)
	fb2.Varchar("id", relNode.Id)
	fb2.Varchar("name", relNode.Name)

	_, err = storeDao.Write(ctx, fb2.Format(fmtStr2), data)

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
	const fmtStr = `
		MATCH (n$<labels>) WHERE n.id=$<id> SET $<props> WITH n 
		MERGE (n)-[:same]->(m$<labels>:same{name:$<name>}) WITH n
		MATCH (n)-[r:same]->(b) WHERE b.name<>$<name> DELETE r
	`
	fb := stringutils.NewFmtBuilder()
	fb.Varchar("name", node.Name)
	fb.Varchar("id", node.Id)
	fb.String("labels", labels)
	fb.String("props", props)

	cypher := fb.Format(fmtStr)

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

	//before := record.BeforeMap()
	after := record.AfterMap()

	node := model2.NewNode(after, record.DBSchema)
	rel := model2.NewRelation(after, record.DBSchema)
	//oldName, _ := maputils.GetString(before, "name", "")

	//tenantId, _ := appctx.GetTenantId(ctx)
	labels := fmt.Sprintf(`:master:case_%s:tenant_%s`, node.CaseId, node.TenantId)

	if isChangedRelType {
		fmtStr := `
			MATCH (n$<labels>)-[r{id:$<id>}]->(m$<labels>) DELETE r 
		`
		fb := stringutils.NewFmtBuilder()
		fb.String("labels", labels)
		fb.String("relType", rel.RelType)
		fb.Varchar("id", node.Id)

		fb.Varchar("tenantId", rel.TenantId)
		fb.Varchar("caseId", rel.CaseId)
		fb.Varchar("sourceType", rel.SourceType)
		fb.Varchar("source", rel.Source)
		fb.Varchar("target", rel.Target)
		fb.Varchar("keywords", rel.Keywords)
		fb.Varchar("description", "")
		fb.Varchar("table", rel.Table)

		if err := d.write(ctx, fmtStr, fb, nil); err != nil {
			panic(err)
		}
		fmtStr = `
			MATCH (n1$<labels>{id:$<source>}), (m1$<labels>{id:$<target>}) WITH n1,m1
			CREATE (n1)-[r1:$<relType>{id:$<id>,case_id:$<caseId>,tenant_id:$<tenantId>,source:$<source>,table:$<table>,
				target:$<target>,source_ids:$<id>,source_type:$<sourceType>,keywords:$<keywords>,description:$<description>
			}]->(m1)
		`
		if err := d.write(ctx, fmtStr, fb, nil); err != nil {
			panic(err)
		}
	}

	if isRename {
		fmtStr := `
		MATCH (n$<labels>{id:$<id>}) SET n.name=$<name> WITH n
		MERGE (m$<labels>:same{name:$<name>}) WITH n, m
		MERGE (n)-[r:same]->(m) WITH n
		MATCH (n)-[rd:same]->(md) WHERE md.name<>$<name> DELETE rd`

		fb := stringutils.NewFmtBuilder()
		fb.String("labels", labels)
		fb.Varchar("name", node.Name)
		fb.Varchar("id", node.Id)
		if err := d.write(ctx, fmtStr, fb, nil); err != nil {
			panic(err)
		}
	}

	if record.IsChangedBusFields() {
		fmtStr := `MATCH (n$<labels>{id:$<id>}) SET n.description=$<description>`
		fb := stringutils.NewFmtBuilder()
		fb.String("labels", labels)
		fb.Varchar("id", node.Id)
		fb.Varchar("description", node.Description)
		if err := d.write(ctx, fmtStr, fb, nil); err != nil {
			panic(err)
		}
	}

	return

}

func (d *NodeDao) write(ctx context.Context, fmtStr string, fb *stringutils.FmtBuilder, params map[string]any) error {
	nodeStore := d.GetStore()
	cypher := fb.Format(fmtStr)
	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		logs.Error(ctx, logs.Fields{"err": err, "cypher": cypher, "params": func() any {
			p, _ := json.Marshal(params)
			return string(p)
		}})

	} else {
		logs.InfoMsg(ctx, cypher)
	}
	return err
}

func (d *NodeDao) DeleteMain(ctx context.Context, record *model2.Record, dbSch *dbschema.DBSchema) {
	nodeStore := d.GetStore()
	before := record.BeforeMap()
	node := model2.NewNode(before, record.DBSchema)
	labels := d.getLabels(node.CaseId, node.TenantId)
	cypher := fmt.Sprintf(`MATCH (n%s{id:$id})-[r]->(m) DETACH DELETE n,r`, labels)
	params := map[string]any{
		"id": node.Id,
	}
	_, err := nodeStore.Write(ctx, cypher, params)
	if err != nil {
		panic(err)
	}
}

func (d *NodeDao) DeleteRelNode(ctx context.Context, record *model2.Record) {
	before := record.BeforeMap()
	rel := model2.NewRelation(before, record.DBSchema)
	labels := d.getLabels(rel.CaseId, rel.TenantId)

	fmtStr := `MATCH (n$<labels>{id:$<id>}) DETACH DELETE n`

	fb := stringutils.NewFmtBuilder()
	fb.String("labels", labels)
	fb.Varchar("id", rel.Id)

	err := d.write(ctx, fmtStr, fb, nil)
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
		Id:          d.GetString(data, "id"),
		Name:        d.GetString(data, "name"),
		CaseId:      d.GetString(data, "case_id"),
		TenantId:    d.GetString(data, "tenant_id"),
		SourceIds:   d.GetString(data, "source_ids"),
		Type:        d.GetString(data, "type"),
		Description: d.GetString(data, "description"),
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
