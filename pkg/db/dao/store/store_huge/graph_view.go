package store_huge

import (
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type GraphView struct {
	*graph.GraphView
	resp *gremlin.PostResponseData
}

func NewGraphView(resp *gremlin.PostResponseData) graph.IGraphView {
	return &GraphView{
		resp:      resp,
		GraphView: graph.NewGraphView(),
	}
}

func NewNode(data map[string]any) *graph.Node {
	if data == nil {
		data = make(map[string]any)
	}
	if data["type"].(string) == "vertex" {
		
	}
	node := &graph.Node{
		Id:    data["id"].(string),
		Label: data["label"].(string),
		Props: data["properties"].(map[string]any),
	}

	return node
}
