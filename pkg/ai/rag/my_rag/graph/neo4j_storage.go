package graph

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Neo4jGraphStorage struct {
	dao   idao.Dao[map[string]any]
	store *store_neo4j.Dao[map[string]any]
}

func NewNeo4jGraphStorage() *Neo4jGraphStorage {
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"master"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
		DBSchema:           dbschema.NewDBSchema("graph", "graph"),
	}
	newDao := dao.NewDao[map[string]any](nodeCfg)
	return &Neo4jGraphStorage{
		dao: newDao,
	}
}

func (d *Neo4jGraphStorage) GetKnowledge(ctx context.Context, tenantId string, caseId string, keys []string, maxDeep int) ([]string, error) {
	contents := []string{}
	graphView := d.FindNodes(ctx, tenantId, caseId, keys, maxDeep)
	nodeMap := make(map[string]*graph.Node)
	for _, nodes := range graphView.Nodes {
		for _, node := range nodes {
			nodeMap[node.Nid] = node
			desc, err := maputils.GetString(node.GetProps(), "description", "")
			if desc != "" && err == nil {
				contents = append(contents, desc)
			}
		}
	}
	for _, edges := range graphView.Edges {
		for _, edge := range edges {
			desc, err := maputils.GetString(edge.GetProps(), "description", "")
			if err != nil {
				continue
			}
			toNode := nodeMap[edge.NTo]
			fromNode := nodeMap[edge.NFrom]
			if toNode != nil && fromNode != nil {
				toName, _ := maputils.GetString(toNode.Props, "name", "")
				fromName, _ := maputils.GetString(fromNode.Props, "name", "")
				relType, _ := maputils.GetString(edge.GetProps(), "relType", "")
				if toName != "" && fromName != "" && relType != "" {
					contents = append(contents, fmt.Sprintf("%s%s%s", toName, relType, fromName))
				}
			}
			if desc != "" {
				contents = append(contents, desc)
			}
		}
	}
	return contents, nil
}

// FindNodes
/*
	MATCH p=(n)-[*..5]-(m) 	WHERE n.name IN ['名称1', '名称2'] RETURN p
*/
func (d *Neo4jGraphStorage) FindNodes(ctx context.Context, tenantId string, caseId string, names []string, maxDeep int) *graph.GraphView {
	namesStr := getNames(names)
	cypher := fmt.Sprintf("MATCH p=(n:tenant_%s:master:case_%s)-[*..%d]-(m) WHERE n.name in [%s] OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m", tenantId, caseId, maxDeep, namesStr)
	logs.InfoMsg(ctx, cypher)
	res, err := d.GetStore().Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.NewGraphView()
}

func (d *Neo4jGraphStorage) GetStore() *store_neo4j.Dao[map[string]any] {
	if d.store == nil {
		iStore := d.dao.GetStore().(any)
		nodeStoreDao, ok := iStore.(*store_neo4j.Dao[map[string]any])
		if !ok {
			panic("neo4j store does not implement neo4j.Dao")
		}
		d.store = nodeStoreDao
	}
	return d.store
}

func getNames(names []string) string {
	str := ""
	count := len(names) - 1
	for i, name := range names {
		str += fmt.Sprintf("'%s'", name)
		if i < count {
			str += ","
		}
	}
	return str
}
