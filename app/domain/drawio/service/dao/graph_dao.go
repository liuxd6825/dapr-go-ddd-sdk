package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"strings"
)

type GraphDao struct {
	idao.Dao[*model.Node]
	store *store_neo4j.Dao[*model.Node]
}

func NewGraphDao() *GraphDao {
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"draw"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	newDao := dao.NewDao[*model.Node](nodeCfg)
	return &GraphDao{
		Dao: newDao,
	}
}

func (d *GraphDao) BatchSave(ctx context.Context, batch *model.SaveBatch, drawId string) {
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

func (d *GraphDao) relationsRemoves(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
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

func (d *GraphDao) nodesCreates(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	i := 0
	for _, item := range batch.Nodes.Creates {
		i++
		n := fmt.Sprintf("n%d", i)
		create := fmt.Sprintf("MERGE ($n%s {id:'%s'})  ON CREATE SET $n.name = '%s' ON MATCH SET $n.id='%s', $n.name='%s' \n",
			d.getItemLabels(tenantId, item, drawId), item.Id, item.Name, item.Id, item.Name)
		create = strings.Replace(create, "$n", n, -1)
		cypher.WriteString(create)
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) nodesUpdates(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	/*
		MATCH (n:human:tenant_test:case_1001:draw_D001:draw)
		WHERE n.id IN ['_JAtuelj7QZR3-mwoWm9-1', '_JAtuelj7QZR3-mwoWm9-2']
		SET n.name = CASE
		    WHEN n.id = '_JAtuelj7QZR3-mwoWm9-1' THEN '张三2'
		    WHEN n.id = '_JAtuelj7QZR3-mwoWm9-2' THEN '李四2'
		    ELSE n.name
		END
	*/
	i := 0

	ids := make([]string, 0)
	names := strings.Builder{}
	for _, item := range batch.Nodes.Updates {
		i++
		ids = append(ids, fmt.Sprintf("'%s'", item.Id))
		names.WriteString(fmt.Sprintf("\n WHEN n.id='%s' THEN '%s'", item.Id, item.Name))
	}
	if i > 0 {
		cypher.WriteString(fmt.Sprintf("\nMATCH (n%s)", d.getDrawLabels(drawId)))
		cypher.WriteString(fmt.Sprintf("\nWHERE n.id IN [%s]", strings.Join(ids, ",")))
		cypher.WriteString("\nSET n.name = CASE")
		cypher.WriteString(names.String())
		cypher.WriteString("ELSE n.name END")
		if cypher.Len() > 0 {
			d.write(ctx, cypher.String())
		}
	}

}

func (d *GraphDao) write(ctx context.Context, cypher string) {
	store := d.GetStore()
	_, err := store.Write(ctx, cypher, nil)
	if err != nil {
		logs.Error(ctx, logs.Fields{"cypher": cypher}, err)
		panic(err)
	}
}

func (d *GraphDao) nodesRemove(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	nodeIds := make([]string, 0)
	for _, item := range batch.Nodes.Removes {
		nodeIds = append(nodeIds, fmt.Sprintf("'%s'", item.Id))
	}
	if len(nodeIds) > 0 {
		cypher.WriteString(fmt.Sprintf("\nMATCH (n:draw_%s) WHERE n.id IN [%s] DETACH DELETE n", drawId, strings.Join(nodeIds, ",")))
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) relationsCreate(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
	// 创建节点
	cypher := strings.Builder{}
	i := 0
	for _, item := range batch.Relations.Creates {
		i++
		a := fmt.Sprintf("a%d", i)
		b := fmt.Sprintf("b%d", i)
		r := fmt.Sprintf("r%d", i)
		create := fmt.Sprintf(
			"\nMATCH ($a:draw_%s{id:'%s'}),($b:draw_%s{id:'%s'}) WITH $a,$b MERGE ($a)-[$r:%s{id:'%s'}]->($b) ON MATCH SET $r.name = '%s' ON CREATE SET $r.id='%s',$r.name='%s' ",
			drawId, item.StartId, drawId, item.EndId, item.RelType, item.Id, item.Name, item.Id, item.Name)
		create = strings.Replace(create, "$a", a, -1)
		create = strings.Replace(create, "$b", b, -1)
		create = strings.Replace(create, "$r", r, -1)
		cypher.WriteString(create)
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) relationsUpdate(ctx context.Context, tenantId string, batch *model.SaveBatch, drawId string) {
	// 更新关系
	cypher := strings.Builder{}
	i := 0
	for _, item := range batch.Relations.Updates {
		r := fmt.Sprintf("r%d", i)
		set := fmt.Sprintf("MATCH (:draw_%s{id:'%s'})-[$r:%s{id:'%s'}]->(:draw_%s{id:'%s'}) WHERE $r.id='%s' SET $r.Name='%s' \n", drawId, item.StartId, item.RelType, item.Id, drawId, item.EndId, item.Id, item.Name)
		set = strings.Replace(set, "$r", r, -1)
		cypher.WriteString(set)
	}
	if cypher.Len() > 0 {
		d.write(ctx, cypher.String())
	}
}

func (d *GraphDao) getItemLabels(tenantId string, item *model.Node, drawId string) string {
	return fmt.Sprintf(":%s:tenant_%s:case_%s:draw_%s:draw", item.Label, tenantId, item.CaseId, drawId)
}

func (d *GraphDao) getDrawLabels(drawId string) string {
	return fmt.Sprintf(":draw_%s:draw", drawId)
}

func (d *GraphDao) GetStore() *store_neo4j.Dao[*model.Node] {
	if d.store == nil {
		iStore := d.Dao.GetStore().(any)
		nodeStoreDao, ok := iStore.(*store_neo4j.Dao[*model.Node])
		if !ok {
			panic("neo4j store does not implement neo4j.Dao")
		}
		d.store = nodeStoreDao
	}
	return d.store
}
