package hugedao

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type MasterNodeDao struct {
	*store_huge.VertexDao[*model.MasterNode]
	labels    []string
	graphMeta *schema.Graph
}

func (d *MasterNodeDao) GetConfig() *idao.DaoConfig {
	//TODO implement me
	panic("implement me")
}

// EdgeLabelMasterRel 统一的 master-mater 关系边 label 名;
// type 作为边属性区分具体关系类型。
const EdgeLabelMasterRel = "master_rel"
const VertexMasterNode = "master_node"

func NewMasterNodeDao(client *hugegraph.CommonClient, graphMeta *schema.Graph) dao2.IMasterNodeDao {
	nodeLabels := []string{VertexMasterNode}
	tableName := "master_node"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.MasterNode{}, tableName)
	cfg := &store_huge.Config[*model.MasterNode]{DBSchema: dbSch, Labels: nodeLabels}
	vertexDao := store_huge.NewVertexDao[*model.MasterNode](client, cfg)
	return &MasterNodeDao{
		graphMeta: graphMeta,
		labels:    nodeLabels,
		VertexDao: vertexDao,
	}
}

func (d *MasterNodeDao) GetLabels() []string {
	return d.labels
}

func (d *MasterNodeDao) GraphMeta() *schema.Graph {
	return d.graphMeta
}

// EnsureSchema 触发 schema 初始化(幂等)
func (d *MasterNodeDao) EnsureSchema(ctx context.Context) error {
	client := d.GetStore().GetDBClient().(*hugegraph.CommonClient)
	return EnsureSchema(ctx, client, d.graphMeta)
}

// CreateMain 等价 Neo4j: CREATE (n:master:case_X:tenant_Y {props}) WITH n
//
//	MERGE (m:master:case_X:tenant_Y:same{name}) WITH n,m
//	MERGE (n)-[:same]->(m)
//
// 拆分为: 1) 创建 master 顶点 (按 id 幂等)  2) 创建 same 节点 (若不存在)  3) 连 same 边
// HugeGraph 1.7 Gremlin 不支持 coalesce(__.unfold(), g.addV(...)), 用 if/else 替代
func (d *MasterNodeDao) CreateMain(ctx context.Context, node *model.MasterNode) error {
	dataMap, err := d.VertexDao.ToMap(node)
	if err != nil {
		return err
	}
	props := d.VertexDao.BuildGremlin(dataMap)
	nLabels := getMasterVertexLabel()
	// 1. 创建 master 顶点 (按 id 幂等 — 同 id 已存在则跳过)
	nodeScript := fmt.Sprintf(
		` if (g.V().hasLabel('%s').has(T.id, '%s').hasNext()) { 0 }  else { g.addV('%s').%s.next(); 0 }`,
		nLabels, escapeGremlin(node.Id),
		nLabels, props,
	)
	if _, err := d.Write(ctx, nodeScript, dataMap); err != nil {
		return err
	}

	sameNodeScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('name','%s').hasNext()) { 0 } else { g.addV('%s').property('name','%s').next(); 0 }",
		VertexLabelSame, escapeGremlin(node.Name),
		VertexLabelSame, escapeGremlin(node.Name),
	)
	if _, err := d.Write(ctx, sameNodeScript, nil); err != nil {
		return err
	}

	sameEdgeScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('id','%s').addE('same').to(__.V().hasLabel('%s').has('name','%s')).property('name','%s').next()",
		nLabels, escapeGremlin(node.Id),
		VertexLabelSame, escapeGremlin(node.Name),
		escapeGremlin(node.Name),
	)
	if _, err := d.Write(ctx, sameEdgeScript, nil); err != nil {
		return err
	}
	return nil
}

// CreateRelNode 等价 Neo4j: CREATE 节点 + MATCH 源 + CREATE 关系 + 关联 same 节点
// HugeGraph 1.7 Gremlin 不支持 coalesce(__.unfold(), g.addV(...)), 拆分为 5 个独立脚本
// relType 改为 EdgeLabel "master_rel" 的 type 属性
func (d *MasterNodeDao) CreateRelNode(ctx context.Context, rel *model.MasterRelation, relNode *model.MasterNode) *model.MasterNode {
	props, data := buildNodeProps(relNode)
	relType := rel.RelType

	// 1. 创建 relNode 顶点
	nodeScript := fmt.Sprintf("g.addV('%s').%s.next()", VertexLabelMaster, props)
	if _, err := d.Write(ctx, nodeScript, data); err != nil {
		panic(err)
	}

	// 2. 创建源节点 (若不存在) — if/else 替代 coalesce
	srcScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('id','%s').hasNext()) { 0 } else { g.addV('%s').property('id','%s').property('case_id','%s').property('tenant_id','%s').next(); 0 }",
		VertexLabelMaster, escapeGremlin(rel.SrcId),
		VertexLabelMaster, escapeGremlin(rel.SrcId),
		escapeGremlin(rel.CaseId), escapeGremlin(rel.TenantId),
	)
	if _, err := d.Write(ctx, srcScript, nil); err != nil {
		panic(err)
	}

	// 3. 创建 master_rel 边,type 作为属性
	edgeScript := fmt.Sprintf(`
		g.V().hasLabel('%s').has('id','%s').addE('%s').to(__.V().hasLabel('%s').has('id','%s')).
			property('id','%s').
			property('type','%s').
			property('case_id','%s').
			property('tenant_id','%s').
			property('description','%s').
			property('rel_source_id','%s').
			property('rel_target_id','%s').
			property('src_id','%s').
			property('src_type','%s').
			property('table','%s').
			property('keywords','%s').next()
	`,
		VertexLabelMaster, escapeGremlin(rel.SrcId),
		EdgeLabelMasterRel,
		VertexLabelMaster, escapeGremlin(relNode.Id),
		escapeGremlin(rel.Id),
		escapeGremlin(relType),
		escapeGremlin(rel.CaseId),
		escapeGremlin(rel.TenantId),
		escapeGremlin(rel.Description),
		escapeGremlin(rel.RelSourceId),
		escapeGremlin(rel.RelTargetId),
		escapeGremlin(rel.Id),
		escapeGremlin(rel.SrcType),
		escapeGremlin(rel.Table),
		escapeGremlin(joinKeywords(rel.Keywords)),
	)
	if _, err := d.Write(ctx, edgeScript, nil); err != nil {
		panic(err)
	}

	// 4. 创建 same 节点 (若不存在)
	sameNodeScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('name','%s').hasNext()) { 0 } else { g.addV('%s').property('name','%s').next(); 0 }",
		VertexLabelSame, escapeGremlin(relNode.Name),
		VertexLabelSame, escapeGremlin(relNode.Name),
	)
	if _, err := d.Write(ctx, sameNodeScript, nil); err != nil {
		panic(err)
	}

	// 5. 连 same 边
	sameEdgeScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('id','%s').addE('same').to(__.V().hasLabel('%s').has('name','%s')).property('name','%s').next()",
		VertexLabelMaster, escapeGremlin(relNode.Id),
		VertexLabelSame, escapeGremlin(relNode.Name),
		escapeGremlin(relNode.Name),
	)
	if _, err := d.Write(ctx, sameEdgeScript, nil); err != nil {
		panic(err)
	}

	return nil
}

// UpdateMain 等价 Neo4j: MATCH (n{...}) SET n.* + 重建 :same 边
// 拆分: 1) 更新属性  2) 创建 same 节点  3)连 same 边  4) 删除旧 same 边
// HugeGraph 1.7 Gremlin 不支持 coalesce(__.unfold(), g.addV(...)), 用 if/else 替代
func (d *MasterNodeDao) UpdateMain(ctx context.Context, node *model.MasterNode) {
	props, data := buildUpdateProps(node)

	setScript := fmt.Sprintf("g.V().hasLabel('%s').has('id','%s').%s.next()",
		VertexLabelMaster, escapeGremlin(node.Id), props)
	if _, err := d.Write(ctx, setScript, data); err != nil {
		panic(err)
	}

	sameNodeScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('name','%s').hasNext()) { 0 } else { g.addV('%s').property('name','%s').next(); 0 }",
		VertexLabelSame, escapeGremlin(node.Name),
		VertexLabelSame, escapeGremlin(node.Name),
	)
	if _, err := d.Write(ctx, sameNodeScript, nil); err != nil {
		panic(err)
	}

	sameEdgeScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('id','%s').addE('same').to(__.V().hasLabel('%s').has('name','%s')).property('name','%s').next()",
		VertexLabelMaster, escapeGremlin(node.Id),
		VertexLabelSame, escapeGremlin(node.Name),
		escapeGremlin(node.Name),
	)
	if _, err := d.Write(ctx, sameEdgeScript, nil); err != nil {
		panic(err)
	}

	dropScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('id','%s').outE('same').has('name', neq('%s')).drop().iterate()",
		VertexLabelMaster,
		escapeGremlin(node.Id), escapeGremlin(node.Name),
	)
	if _, err := d.Write(ctx, dropScript, nil); err != nil {
		panic(err)
	}
}

// UpdateRelNode CDC 增量更新,处理 rename / relType 变更 / 业务字段更新
func (d *MasterNodeDao) UpdateRelNode(ctx context.Context, record *dbevent.CDCRecord) error {
	isRename := d.IsRename(record)
	isChangedRelType := d.IsChangedRelType(record)

	after := record.AfterMap()
	node := model.NewMasterNode(after, record.DBSchema)
	rel := model.NewMasterRelation(after, record.DBSchema)

	if isChangedRelType {
		// type 是属性,变更即直接 update 边的 type 属性
		updScript := fmt.Sprintf(
			"g.E().hasLabel('%s').has('id','%s').property('rel_type','%s').iterate()",
			EdgeLabelMasterRel,
			escapeGremlin(rel.Id),
			escapeGremlin(rel.RelType),
		)
		if _, err := d.Write(ctx, updScript, nil); err != nil {
			return err
		}
	}

	if isRename {
		// 1. 更新 name 属性
		setNameScript := fmt.Sprintf(
			"g.V().hasLabel('%s').has('id','%s').property('name','%s').next()",
			VertexLabelMaster,
			escapeGremlin(node.Id), escapeGremlin(node.Name),
		)
		if _, err := d.Write(ctx, setNameScript, nil); err != nil {
			panic(err)
		}
		// 2. 创建 same 节点 (若不存在) — if/else 替代 coalesce
		sameNodeScript := fmt.Sprintf(
			"if (g.V().hasLabel('%s').has('name','%s').hasNext()) { 0 } else { g.addV('%s').property('name','%s').next(); 0 }",
			VertexLabelSame, escapeGremlin(node.Name),
			VertexLabelSame, escapeGremlin(node.Name),
		)
		if _, err := d.Write(ctx, sameNodeScript, nil); err != nil {
			panic(err)
		}
		// 3. 连 same 边
		sameEdgeScript := fmt.Sprintf(
			"g.V().hasLabel('%s').has('id','%s').addE('same').to(__.V().hasLabel('%s').has('name','%s')).property('name','%s').next()",
			VertexLabelMaster, escapeGremlin(node.Id),
			VertexLabelSame, escapeGremlin(node.Name),
			escapeGremlin(node.Name),
		)
		if _, err := d.Write(ctx, sameEdgeScript, nil); err != nil {
			panic(err)
		}
		// 4. 删除指向其它 name 的过期 same 边
		dropScript := fmt.Sprintf(
			"g.V().hasLabel('%s').has('id','%s').outE('same').has('name', neq('%s')).drop().iterate()",
			VertexLabelMaster,
			escapeGremlin(node.Id), escapeGremlin(node.Name),
		)
		if _, err := d.Write(ctx, dropScript, nil); err != nil {
			panic(err)
		}
	}

	if d.isChangedBusFields(record) {
		busScript := fmt.Sprintf(
			"g.V().hasLabel('%s').has('id','%s').property('description','%s').next()",
			VertexLabelMaster,
			escapeGremlin(node.Id), escapeGremlin(node.Description),
		)
		if _, err := d.Write(ctx, busScript, nil); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMain 等价 Neo4j: MATCH (n{id})-[r]->(m) DETACH DELETE n,r
func (d *MasterNodeDao) DeleteMain(ctx context.Context, record *dbevent.CDCRecord, _ *store2.DBSchema) {
	before := record.BeforeMap()
	node := model.NewMasterNode(before, record.DBSchema)
	script := fmt.Sprintf("g.V().hasLabel('%s').has('id','%s').drop().iterate()",
		VertexLabelMaster, escapeGremlin(node.Id))
	if _, err := d.Write(ctx, script, nil); err != nil {
		panic(err)
	}
}

// DeleteRelNode 删除关系节点
func (d *MasterNodeDao) DeleteRelNode(ctx context.Context, record *dbevent.CDCRecord) {
	before := record.BeforeMap()
	rel := model.NewMasterRelation(before, record.DBSchema)
	script := fmt.Sprintf("g.V().hasLabel('%s').has('id','%s').drop().iterate()",
		VertexLabelMaster, escapeGremlin(rel.Id))
	if _, err := d.Write(ctx, script, nil); err != nil {
		panic(err)
	}
}

// ClearAll 清空所有节点和边
func (d *MasterNodeDao) ClearAll(ctx context.Context) {
	script := "g.V().drop().iterate();g.E().drop().iterate()"
	if _, err := d.Write(ctx, script, nil); err != nil {
		panic(err)
	}
}

// NewNode 从 map 构造 MasterNode
func (d *MasterNodeDao) NewNode(data map[string]any) *model.MasterNode {
	return &model.MasterNode{
		Id:          d.GetString(data, "id"),
		Name:        d.GetString(data, "name"),
		CaseId:      d.GetString(data, "case_id"),
		TenantId:    d.GetString(data, "tenant_id"),
		SrcId:       d.GetString(data, "src_id"),
		Type:        d.GetString(data, "type"),
		Description: d.GetString(data, "description"),
	}
}

func (d *MasterNodeDao) GetString(vMap map[string]any, dbName string) string {
	val, _ := maputils.GetString(vMap, dbName, "")
	return val
}

func (d *MasterNodeDao) GetInt64(vMap map[string]any, dbName string) int64 {
	val, _ := maputils.GetInt64(vMap, dbName, -1)
	return val
}

// IsRename 是否数据更新
func (d *MasterNodeDao) IsRename(r *dbevent.CDCRecord) bool {
	if r.OpType == dbevent.OpTypeUpdate {
		name := d.graphMeta.Name
		if name == "" {
			name = "name"
		}
		newName, _ := maputils.GetString(r.After, name, "")
		oldName, _ := maputils.GetString(r.Before, name, "")
		if newName != oldName {
			return true
		}
	}
	return false
}

// IsChangedRelType 是否数据更新
func (d *MasterNodeDao) IsChangedRelType(r *dbevent.CDCRecord) bool {
	graphMeta := d.graphMeta
	if graphMeta.IsRelTypeField() {
		if r.OpType == "u" {
			refType := graphMeta.GetRelTypeField()
			newType, _ := maputils.GetString(r.After, refType, "")
			oldType, _ := maputils.GetString(r.Before, refType, "")
			if newType != oldType {
				return true
			}
		}
	}
	return false
}

// isChangedBusFields 是否更新业务字段
func (d *MasterNodeDao) isChangedBusFields(record *dbevent.CDCRecord) bool {
	graphMeta := d.graphMeta
	if graphMeta.IsRelTypeField() {
		refType := graphMeta.GetRelTypeField()
		if record.OpType == dbevent.OpTypeUpdate {
			for k := range record.After {
				if k != refType {
					return true
				}
			}
		}
		return false
	}
	return true
}

func (d *MasterNodeDao) GetStore() store2.IStore[*model.MasterNode] {
	return nil
}

// ---- helpers ----

func getMasterVertexLabel() string {
	return VertexLabelMaster
}

func buildNodeProps(node *model.MasterNode) (string, map[string]any) {
	nodeType := node.Type
	if nodeType == "" {
		nodeType = "master"
	}

	data := map[string]any{
		"id":          node.Id,
		"name":        node.Name,
		"case_id":     node.CaseId,
		"tenant_id":   node.TenantId,
		"keywords":    node.Keywords,
		"src_id":      node.SrcId,
		"src_url":     node.SrcUrl,
		"src_type":    node.SrcType,
		"src_name":    "",
		"type":        nodeType,
		"node_type":   nodeType,
		"description": node.Description,
		"table":       node.Table,
		"tags":        "",
	}
	sb := strings.Builder{}
	count := len(data)
	i := 0
	for k, v := range data {
		i++
		key := k
		if key == "id" {
			key = "T.id"
		} else {
			key = "'" + key + "'"
		}
		val := v.(string)
		sb.WriteString(fmt.Sprintf("property(%s,'%s')", key, escapeGremlin(val)))
		if i < count {
			sb.WriteString(".")
		}
	}
	/*	props := fmt.Sprintf(
		"property(T.id,'%s').property('name','%s').property('case_id','%s').property('tenant_id','%s').property('src_id','%s').property('src_type','%s').property('type','%s').property('node_type','%s').property('description','%s').property('table','%s')",
		escapeGremlin(node.Id), escapeGremlin(node.Name), escapeGremlin(node.CaseId),
		escapeGremlin(node.TenantId), escapeGremlin(node.SrcId), escapeGremlin(node.SrcType),
		escapeGremlin(nodeType), escapeGremlin(nodeType),
		escapeGremlin(node.Description), escapeGremlin(node.Table),
	)*/
	//return props, data
	return sb.String(), data
}

func buildUpdateProps(node *model.MasterNode) (string, map[string]any) {
	v := reflect.ValueOf(node).Elem()
	t := v.Type()
	data := map[string]any{}
	props := ""
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := f.Tag.Get("gorm")
		if name == "" {
			continue
		}
		// 提取 column:"xxx"
		dbName := extractGormColumn(name)
		if dbName == "" || dbName == "id" {
			continue
		}
		val := fmt.Sprintf("%v", v.Field(i).Interface())
		data[dbName] = val
		props += fmt.Sprintf(".property('%s','%s')", dbName, escapeGremlin(val))
	}
	return props, data
}

func extractGormColumn(tag string) string {
	parts := splitTag(tag)
	for _, p := range parts {
		if len(p) > 7 && p[:7] == "column:" {
			return p[7:]
		}
	}
	return ""
}

func splitTag(tag string) []string {
	out := []string{}
	cur := ""
	for _, r := range tag {
		if r == ';' {
			if cur != "" {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func escapeGremlin(s string) string {
	out := ""
	for _, r := range s {
		switch r {
		case '\'':
			out += "\\'"
		case '\\':
			out += "\\\\"
		default:
			out += string(r)
		}
	}
	return out
}

func joinKeywords(kw []string) string {
	if len(kw) == 0 {
		return ""
	}
	out := ""
	for i, k := range kw {
		if i > 0 {
			out += ","
		}
		out += escapeGremlin(k)
	}
	return out
}

// 防止编译器报告 stringutils/maputils 未使用
var (
	_ = stringutils.NewFmtBuilder
	_ = maputils.GetInt64
)
