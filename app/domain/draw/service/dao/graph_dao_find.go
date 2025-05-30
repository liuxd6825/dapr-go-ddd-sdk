package dao

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

func (d *GraphDao) FindGraphByDrawId(ctx context.Context, drawId string) *graph.GraphView {
	tenantId := appctx.GetTenantId2(ctx)
	cypher := fmt.Sprintf("MATCH (n:tenant_%s:draw_%s) OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, drawId)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *GraphDao) AddNodeList(target map[string]*model.NodeView, source []any) {
	for _, item := range source {
		if dbNode, ok := item.(dbtype.Node); ok {
			node := model.NewNodeView()
			node.Id, _ = maputils.GetString(dbNode.Props, "id", "")
			node.Name, _ = maputils.GetString(dbNode.Props, "name", "")
			node.Nid = dbNode.ElementId
			node.Tags = dbNode.Labels
			node.Props = dbNode.Props
			target[node.Id] = node
		}
	}
}

func (d *GraphDao) AddRelationList(target map[string]*model.RelationView, source []any) {
	for _, item := range source {
		if dbRel, ok := item.(dbtype.Relationship); ok {
			rel := model.NewRelationView()
			rel.NId = dbRel.ElementId
			rel.Id, _ = maputils.GetString(dbRel.Props, "id", "")
			rel.RelType = dbRel.Type
			rel.Props = dbRel.Props
			rel.StartId = dbRel.StartElementId
			rel.EndId = dbRel.EndElementId
			target[rel.Id] = rel
		}
	}
}
