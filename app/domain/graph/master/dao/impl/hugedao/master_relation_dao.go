package hugedao

import (
	"context"
	"fmt"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	store_huge "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

type BusRelationDao struct {
	idao.Dao[*model2.MasterRelation]
	nodeDao dao2.IMasterNodeDao
	*Base[*model2.MasterRelation]
	graphMeta *schema.Graph
}

func NewBusRelationDao(dbSch *store.DBSchema, nodeDao dao2.IMasterNodeDao, graphMeta *schema.Graph) *BusRelationDao {
	relCfg := &dao.DaoConfig{
		DBKey:              DB_KEY_HUGE,
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

// EnsureSchema 触发 schema 初始化
func (d *BusRelationDao) EnsureSchema(ctx context.Context) error {
	client := d.GetStore().GetDBClient().(*hugegraph.CommonClient)
	return EnsureSchema(ctx, client, d.graphMeta)
}

// CreateSameName 占位: 与 Neo4j 版一致,此处保留为 no-op
func (d *BusRelationDao) CreateSameName(_ context.Context, _ *model2.MasterNode) {
}

// FindByName 通过 Gremlin 直接查询 master 节点上 name 匹配的第一条
func (d *BusRelationDao) FindByName(ctx context.Context, name string) *model2.MasterRelation {
	script := fmt.Sprintf(
		"g.V().hasLabel('%s').has('name','%s').valueMap(true)",
		store_huge.VertexLabelMaster, escapeGremlin(name),
	)
	res, err := d.GetStore().Gremlin(ctx, script, nil)
	if err != nil {
		panic(err)
	}
	records := store_huge.NewHugeResult(res).Records()
	if len(records) == 0 {
		return nil
	}
	return mapToMasterRelation(records[0])
}

// UpdateRecord 等价 Neo4j: MERGE 源节点 + MERGE 目标节点 + MERGE 关系
// rel_type 改为 master_rel 边的属性
func (d *BusRelationDao) UpdateRecord(ctx context.Context, record *restapi.CDCRecord) *model2.MasterNode {
	after := record.AfterMap()
	node := model2.NewMasterNode(after, record.DBSchema)
	rel := model2.NewMasterRelation(after, record.DBSchema)
	nodeName, _ := after["name"].(string)

	tenantId, _ := ctx.Value("tenantId").(string)
	if tenantId == "" {
		tenantId = rel.TenantId
	}

	relType := rel.RelType

	// 拆分为 4 个独立脚本 (HugeGraph 1.7 不支持 fold+coalesce+addV)
	// 1. 创建 src 顶点 (若不存在)
	srcScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('id','%s').hasNext()) { 0 } else { g.addV('%s').property('id','%s').property('name','%s').next(); 0 }",
		store_huge.VertexLabelMaster, escapeGremlin(rel.RelSourceId),
		store_huge.VertexLabelMaster,
		escapeGremlin(node.Id), escapeGremlin(node.Name),
	)
	if _, err := d.GetStore().Gremlin(ctx, srcScript, nil); err != nil {
		panic(err)
	}

	// 2. 创建 dst 顶点 (若不存在)
	dstScript := fmt.Sprintf(
		"if (g.V().hasLabel('%s').has('name','%s').hasNext()) { 0 } else { g.addV('%s').property('tenant_id','%s').property('name','%s').property('description','%s').next(); 0 }",
		store_huge.VertexLabelMaster, escapeGremlin(nodeName),
		store_huge.VertexLabelMaster,
		escapeGremlin(tenantId), escapeGremlin(nodeName), escapeGremlin(rel.Description),
	)
	if _, err := d.GetStore().Gremlin(ctx, dstScript, nil); err != nil {
		panic(err)
	}

	// 3. 创建或获取 master_rel 边 (若已存在则跳过)
	relExistScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('id','%s').outE('%s').has('id','%s').count()",
		store_huge.VertexLabelMaster, escapeGremlin(rel.RelSourceId),
		EdgeLabelMasterRel, escapeGremlin(rel.Id),
	)
	relE := d.GetStore()
	res, err := relE.Gremlin(ctx, relExistScript, nil)
	if err != nil {
		panic(err)
	}
	count := 0
	if list, ok := res.Result.Data.([]any); ok && len(list) > 0 {
		if v, ok := list[0].(float64); ok {
			count = int(v)
		}
	}
	if count == 0 {
		creScript := fmt.Sprintf(
			"g.V().hasLabel('%s').has('id','%s').addE('%s').to(__.V().hasLabel('%s').has('name','%s')).property('id','%s').property('rel_type','%s').property('tenant_id','%s').property('case_id','%s').property('description','%s').next()",
			store_huge.VertexLabelMaster, escapeGremlin(rel.RelSourceId),
			EdgeLabelMasterRel,
			store_huge.VertexLabelMaster, escapeGremlin(nodeName),
			escapeGremlin(rel.Id),
			escapeGremlin(relType),
			escapeGremlin(tenantId),
			escapeGremlin(rel.CaseId),
			escapeGremlin(rel.Description),
		)
		if _, err := d.GetStore().Gremlin(ctx, creScript, nil); err != nil {
			panic(err)
		}
	}

	return nil
}

// DeleteRecord 删除关系,若端点仅剩此关系则一并删除端点
func (d *BusRelationDao) DeleteRecord(ctx context.Context, record *restapi.CDCRecord) {
	after := record.AfterMap()
	rel := model2.NewMasterRelation(after, record.DBSchema)
	name, _ := after["name"].(string)

	// 1. 删除关系
	delRelScript := fmt.Sprintf(
		"g.E().hasLabel('%s').has('id','%s').drop().iterate()",
		EdgeLabelMasterRel, escapeGremlin(rel.Id),
	)
	if _, err := d.GetStore().Gremlin(ctx, delRelScript, nil); err != nil {
		panic(err)
	}

	// 2. 清理孤立端点
	delOrphanScript := fmt.Sprintf(
		"g.V().hasLabel('%s').has('name','%s').where(__.inE().count().is(eq(0)).and().outE().count().is(eq(0))).drop().iterate()",
		store_huge.VertexLabelMaster, escapeGremlin(name),
	)
	if _, err := d.GetStore().Gremlin(ctx, delOrphanScript, nil); err != nil {
		panic(err)
	}
}

func mapToMasterRelation(m map[string]any) *model2.MasterRelation {
	get := func(k string) string {
		if v, ok := m[k]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
	return &model2.MasterRelation{
		Id:          get("id"),
		CaseId:      get("case_id"),
		TenantId:    get("tenant_id"),
		RelSourceId: get("rel_source_id"),
		RelTargetId: get("rel_target_id"),
		SourceId:    get("source_id"),
		SourceType:  get("source_type"),
		Description: get("description"),
		Table:       get("table"),
	}
}

func (d *BusRelationDao) GetStore() *store_huge.Dao[*model2.MasterRelation] {
	iStore := d.Dao.GetStore().(any)
	storeDao, ok := iStore.(*store_huge.Dao[*model2.MasterRelation])
	if !ok {
		panic("hugegraph store does not implement store_huge.Dao")
	}
	return storeDao
}
