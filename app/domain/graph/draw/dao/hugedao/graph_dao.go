package hugedao

import (
	"context"
	"fmt"
	"strings"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/hugegraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type GraphDao struct {
	store *hugegraph.Store
}

type Prop string

/*
fb.String("label", d.getItemLabels(tenantId, item, drawId))
fb.Varchar("id", item.Id)
fb.Varchar("name", item.Name)
fb.Varchar("desc", item.Description)
fb.Varchar("type", item.Type)
fb.Varchar("case_id", item.CaseId)
fb.Varchar("draw_id", item.DrawId)
fb.Varchar("source_ids", item.SourceIds)
fb.Varchar("source_type", item.SourceType)
fb.Varchar("source_url", item.SourceUrl)
fb.Varchar("table", item.Table)
*/
const (
	PropLabel      = "label"
	PropId         = "id"
	PropName       = "name"
	PropDesc       = "desc"
	PropType       = "type"
	PropCaseId     = "case_id"
	PropDrawId     = "draw_id"
	PropSourceId   = "source_id"
	PropSourceType = "source_type"
	PropSourceUrl  = "source_url"
	PropTable      = "table"
)

func NewGraphDao(store *hugegraph.Store) *GraphDao {

	return &GraphDao{
		store: store,
	}
}

func (d *GraphDao) BatchSave(ctx context.Context, batch *model2.SaveBatch, drawId string) {
	tenantId := appctx.GetTenantId2(ctx)
	sb := &strings.Builder{}
	// 删除关系
	d.relationsRemoves(ctx, sb, tenantId, batch, drawId)

	// 创建节点
	d.nodesCreates(ctx, sb, tenantId, batch, drawId)

	// 更新节点
	d.nodesUpdates(ctx, sb, tenantId, batch, drawId)

	// 删除节点
	d.nodesRemove(ctx, sb, tenantId, batch, drawId)

	// 创建关系
	d.relationsCreate(ctx, sb, tenantId, batch, drawId)

	// 更新关系
	d.relationsUpdate(ctx, sb, tenantId, batch, drawId)

	d.write(ctx, sb.String())

}

func (d *GraphDao) relationsRemoves(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 删除关系
	relIds := make([]string, 0)
	for _, item := range batch.Relations.Removes {
		relIds = append(relIds, fmt.Sprintf("'%s'", item.Id))
	}
	if len(relIds) > 0 {
		cypher := fmt.Sprintf("g.E(%s).drop(); \n", strings.Join(relIds, ","))
		sb.WriteString(cypher)
	}
}

func (d *GraphDao) nodesCreates(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	for _, item := range batch.Nodes.Creates {
		create := `
	g.V($<id>).fold().coalesce(
    unfold()
        .property('id', $<id>)
        .property('name', $<name>)
        .property('type', $<type>)
        .property('case_id', $<case_id>)
        .property('description', $<desc>)
        .property('draw_id', $<draw_id>)
        .property('source_ids', $<source_ids>)
        .property('source_type', $<source_type>)
        .property('source_url', $<source_url>)
        .property('table', $<table>),
    addV($<label>)
        .property(T.id, $<id>)  
        .property('id', $<id>)
        .property('name', $<name>)
); 
`
		fb := stringutils.NewFmtBuilder()
		fb.String(PropLabel, d.getItemLabels(tenantId, item, drawId))
		fb.Varchar(PropId, item.Id)
		fb.Varchar(PropName, item.Name)
		fb.Varchar(PropDesc, item.Description)
		fb.Varchar(PropType, item.Type)
		fb.Varchar(PropCaseId, item.CaseId)
		fb.Varchar(PropDrawId, item.DrawId)
		fb.Varchar(PropSourceIds, item.SrcId)
		fb.Varchar(PropSourceType, item.SrcType)
		fb.Varchar(PropSourceUrl, item.SrcUrl)
		fb.Varchar(PropTable, item.Table)
		cypher := fb.Format(create)
		sb.WriteString(cypher + "\n")
	}
}

func (d *GraphDao) nodesUpdates(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	/*
		MATCH (n:human:tenant_test:case_1001:draw_D001:draw) WHERE n.id='' SET n.name='', n.source='', n.target=‘’
	*/
	for _, item := range batch.Nodes.Updates {
		update := `
     g.V($<id>)
    .property('name', $<name>)
    .property('type', $<type>)
    .property('case_id', $<case_id>)
    .property('description', $<desc>)
    .property('draw_id', $<draw_id>)
    .property('src_id', $<source_id>)
    .property('src_type', $<source_type>)
    .property('source_url', $<source_url>)
    .property('table', $<table>)
`
		fb := stringutils.NewFmtBuilder()
		fb.String("n", "n")
		fb.Varchar(PropId, item.Id)
		fb.Varchar(PropName, item.Name)
		fb.Varchar(PropDesc, item.Description)
		fb.Varchar(PropType, item.Type)
		fb.Varchar(PropCaseId, item.CaseId)
		fb.Varchar(PropDrawId, item.DrawId)
		fb.Varchar(PropSourceIds, item.SrcId)
		fb.Varchar(PropSourceType, item.SrcType)
		fb.Varchar(PropSourceUrl, item.SrcUrl)
		fb.Varchar(PropTable, item.Table)

		cypher := fb.Format(update)
		sb.WriteString(cypher + "\n")
	}

}

func (d *GraphDao) write(ctx context.Context, cypher string) {
	logs.Info(ctx, logs.Fields{"cypher": cypher})
	store := d.store
	_, err := store.Write(ctx, cypher, nil)
	if err != nil {
		logs.Error(ctx, logs.Fields{"cypher": cypher}, err)
		panic(err)
	}
}

func (d *GraphDao) nodesRemove(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	nodeIds := make([]string, 0)
	for _, item := range batch.Nodes.Removes {
		nodeIds = append(nodeIds, fmt.Sprintf("'%s'", item.Id))
	}
	if len(nodeIds) > 0 {
		cypherStr := " g.V($<nodeIds>).drop(); \n"

		fb := stringutils.NewFmtBuilder()
		fb.String("drawId", drawId)
		fb.StringsJoin("nodeIds", nodeIds, ",")
		cypher.WriteString(fb.Format(cypherStr))
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) relationsCreate(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	for _, item := range batch.Relations.Creates {
		create := `
g.V($<source>).addE($<label>).to(g.V($<target>))
 .property('name', $<name>)
 .property('source', $<source>)
 .property('target', $<target>)
 .property('relType', $<relType>)
 .property('caseId', $<caseId>)
 .property('drawId', $<drawId>)
 .property('desc', $<desc>) \n
`
		fb := stringutils.NewFmtBuilder()
		fb.Varchar("id", item.Id)
		fb.String("caseId", item.CaseId)
		fb.String("drawId", drawId)
		fb.String("label", "draw_rel")
		fb.Varchar("source", item.SourceId)
		fb.Varchar("target", item.TargetId)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		fb.Varchar("relType", item.RelType)
		cypher := fb.Format(create)
		d.write(ctx, cypher)
	}
}

func (d *GraphDao) relationsUpdate(ctx context.Context, sb *strings.Builder, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 更新关系
	for _, item := range batch.Relations.Updates {
		set := `
g.V().hasLabel('draw').has('id', $<source>)
    .outE($<relType>).has('id', $<id>)
    .where(inV().hasLabel('draw').has('id', $<target>))
    .property('name', $<name>)
    .property('source', $<source>)
    .property('target', $<target>)
    .property('description', $<desc>)
    .property('keywords', $<keywords>)
    .property('source_type', 'draw')
    .property('draw_id', $<drawId>);\n
`
		fb := stringutils.NewFmtBuilder()
		fb.String("r", "r")
		fb.String("caseId", item.CaseId)
		fb.String("drawId", drawId)
		fb.Varchar("id", item.Id)
		fb.Varchar("source", item.Source)
		fb.Varchar("target", item.Target)
		fb.String("relType", "draw")
		//fb.Varchar("relType", item.RelType)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		fb.Varchar("keywords", item.RelType)
		fb.Varchar("source_ids", batch.DrawName)
		sb.WriteString(fb.Format(set))
	}
}

func (d *GraphDao) getItemLabels(tenantId string, item *model2.Node, drawId string) string {
	return "draw"
	/*	if item.Type != "" {
			return fmt.Sprintf(":%s:tenant_%s:case_%s:draw_%s:draw", item.Type, tenantId, item.CaseId, drawId)
		}
		return fmt.Sprintf(":tenant_%s:case_%s:draw_%s:draw", tenantId, item.CaseId, drawId)*/
}

func (d *GraphDao) getDrawLabels(drawId string) string {
	return fmt.Sprintf(":draw_%s:draw", drawId)
}
