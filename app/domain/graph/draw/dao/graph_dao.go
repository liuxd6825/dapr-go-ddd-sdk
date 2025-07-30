package dao

import (
	"context"
	"fmt"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"strings"
)

type GraphDao struct {
	idao.Dao[*model2.Node]
	store *store_neo4j.Dao[*model2.Node]
}

func NewGraphDao(neo4jDBKey string) *GraphDao {
	nodeCfg := &dao.NewConfig{
		DBKey:              neo4jDBKey,
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"draw"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	newDao := dao.NewDao[*model2.Node](nodeCfg)
	return &GraphDao{
		Dao: newDao,
	}
}

func (d *GraphDao) BatchSave(ctx context.Context, batch *model2.SaveBatch, drawId string) {
	tenantId := appctx.GetTenantId2(ctx)

	// 删除关系
	d.relationsRemoves(ctx, tenantId, batch, drawId)

	// 创建节点
	d.nodesCreates(ctx, tenantId, batch, drawId)

	// 更新节点
	d.nodesUpdates(ctx, tenantId, batch, drawId)

	// 删除节点
	d.nodesRemove(ctx, tenantId, batch, drawId)

	// 创建关系
	d.relationsCreate(ctx, tenantId, batch, drawId)

	// 更新关系
	d.relationsUpdate(ctx, tenantId, batch, drawId)

}

func (d *GraphDao) relationsRemoves(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 删除关系
	relIds := make([]string, 0)
	for _, item := range batch.Relations.Removes {
		relIds = append(relIds, fmt.Sprintf("\"%s\"", item.Id))
	}
	if len(relIds) > 0 {
		cypher := fmt.Sprintf("MATCH ()-[r1]->() WHERE r1.id IN [%s] DELETE r1 \n", strings.Join(relIds, ","))
		d.write(ctx, cypher)
	}
}

func (d *GraphDao) nodesCreates(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	for _, item := range batch.Nodes.Creates {
		create := "MERGE ($<n>$<label>{id:$<id>}) ON CREATE SET $<n>.name=$<name> ON MATCH SET $<n>.id=$<id>,$<n>.name=$<name>,$<n>.description=$<desc>; \n"
		fb := stringutils.NewFmtBuilder()
		fb.String("n", fmt.Sprintf("n%d", 1))
		fb.String("caseId", item.CaseId)
		fb.String("drawId", drawId)
		fb.String("label", d.getItemLabels(tenantId, item, drawId))
		fb.Varchar("id", item.Id)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		fb.Varchar("type", item.Type)
		d.write(ctx, fb.Format(create))
	}
}

func (d *GraphDao) nodesUpdates(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	/*
		MATCH (n:human:tenant_test:case_1001:draw_D001:draw) WHERE n.id='' SET n.name='', n.source='', n.target=‘’
	*/
	for _, item := range batch.Nodes.Updates {
		update := "MATCH ($<n>$<labels>{id:$<id>}) SET $<n>.name=$<name>,$<n>.description=$<desc> \n"
		fb := stringutils.NewFmtBuilder()
		fb.String("n", "n")
		fb.String("labels", d.getDrawLabels(drawId))
		fb.Varchar("id", item.Id)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		d.write(ctx, fb.Format(update))
	}

}

func (d *GraphDao) write(ctx context.Context, cypher string) {
	logs.Info(ctx, logs.Fields{"cypher": cypher})
	store := d.GetStore()
	_, err := store.Write(ctx, cypher, nil)
	if err != nil {
		logs.Error(ctx, logs.Fields{"cypher": cypher}, err)
		panic(err)
	}
}

func (d *GraphDao) nodesRemove(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	nodeIds := make([]string, 0)
	for _, item := range batch.Nodes.Removes {
		nodeIds = append(nodeIds, fmt.Sprintf("'%s'", item.Id))
	}
	if len(nodeIds) > 0 {
		cypherStr := " MATCH (n:draw_$<drawId>) WHERE n.id IN [$<nodeIds>] DETACH DELETE n  \n"
		fb := stringutils.NewFmtBuilder()
		fb.String("drawId", drawId)
		fb.StringsJoin("nodeIds", nodeIds, ",")
		cypher.WriteString(fb.Format(cypherStr))
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) relationsCreate(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 创建节点
	for _, item := range batch.Relations.Creates {
		create := `
		MATCH ($<a>:draw_$<drawId>{id:$<source>}),($<b>:draw_$<drawId>{id:$<target>})
		WITH $<a>,$<b> MERGE ($<a>)-[$<r>:$<relType>{id:$<id>}]->($<b>) 
		ON MATCH  SET $<r>.name=$<name>,$<r>.source=$<source>,$<r>.target=$<target>,$<r>.description=$<r>.description+$<desc>,$<r>.keywords=$<keywords>
		ON CREATE SET $<r>.id=$<id>,$<r>.name=$<name>,$<r>.source=$<source>,$<r>.target=$<target>,$<r>.description=$<desc>,$<r>.keywords=$<keywords>
`
		fb := stringutils.NewFmtBuilder()
		fb.String("a", "a")
		fb.String("b", "b")
		fb.String("r", "r")
		fb.String("caseId", item.CaseId)
		fb.String("drawId", drawId)
		fb.String("relType", item.RelType)

		fb.Varchar("id", item.Id)
		fb.Varchar("source", item.Source)
		fb.Varchar("target", item.Target)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		fb.Varchar("keywords", item.RelType)

		d.write(ctx, fb.Format(create))
	}
}

func (d *GraphDao) relationsUpdate(ctx context.Context, tenantId string, batch *model2.SaveBatch, drawId string) {
	// 更新关系
	for _, item := range batch.Relations.Updates {
		set := `
		MATCH (:draw_$<drawId>{id:$<source>})-[$<r>:$<relType>{id:$<id>}]->(:draw_$<drawId>{id:$<target>}) 
		WHERE $<r>.id=$<id> 
		SET $<r>.name=$<name>,$<r>.source=$<source>,$<r>.target=$<target>,$<r>.description=$<desc>,$<r>.keywords=$<keywords>,$<r>.source_type='draw' 
`
		fb := stringutils.NewFmtBuilder()
		fb.String("r", "r")
		fb.String("caseId", item.CaseId)
		fb.String("drawId", drawId)
		fb.Varchar("id", item.Id)
		fb.Varchar("source", item.Source)
		fb.Varchar("target", item.Target)
		fb.Varchar("relType", item.RelType)
		fb.Varchar("name", item.Name)
		fb.Varchar("desc", item.Description)
		fb.Varchar("keywords", item.RelType)
		fb.Varchar("source_ids", batch.DrawName)
		d.write(ctx, fb.Format(set))
	}
}

func (d *GraphDao) getItemLabels(tenantId string, item *model2.Node, drawId string) string {
	if item.Type != "" {
		return fmt.Sprintf(":%s:tenant_%s:case_%s:draw_%s:draw", item.Type, tenantId, item.CaseId, drawId)
	}
	return fmt.Sprintf(":tenant_%s:case_%s:draw_%s:draw", tenantId, item.CaseId, drawId)
}

func (d *GraphDao) getDrawLabels(drawId string) string {
	return fmt.Sprintf(":draw_%s:draw", drawId)
}

func (d *GraphDao) GetStore() *store_neo4j.Dao[*model2.Node] {
	if d.store == nil {
		iStore := d.Dao.GetStore().(any)
		nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model2.Node])
		if !ok {
			panic("neo4j store does not implement neo4j.Dao")
		}
		d.store = nodeStoreDao
	}
	return d.store
}
