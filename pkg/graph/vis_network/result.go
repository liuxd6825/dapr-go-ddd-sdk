package vis_network

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"

type NodeMap map[string]*Node
type EdgeMap map[string]*Edge

// Result Vis数据
type Result struct {
	Nodes     map[string]NodeMap `json:"nodes"`
	Edges     map[string]EdgeMap `json:"edges"`
	Heads     map[string]any     `json:"heads"`
	PageTotal *int64             `json:"pageTotal"`
	PageSize  *int64             `json:"pageSize"`
	PageNum   *int64             `json:"pageNum"`
	Total     *int64             `json:"total"`
}

func NewResult() *Result {
	return &Result{
		Heads: nil,
		Nodes: map[string]NodeMap{},
		Edges: map[string]EdgeMap{},
	}
}

func NewResultWithGraphView(graphView *graph.GraphView) *Result {
	result := &Result{
		Heads: nil,
		Nodes: map[string]NodeMap{},
		Edges: map[string]EdgeMap{},
	}
	for dataKey, graphNodes := range graphView.Nodes {
		nodes := NodeMap{}
		for key, graphNode := range graphNodes {
			nodes[key] = &Node{
				Node: *graphNode,
			}
		}
		result.Nodes[dataKey] = nodes
	}
	for dataKey, graphEdges := range graphView.Edges {
		edges := EdgeMap{}
		for key, graphEdge := range graphEdges {
			edges[key] = &Edge{
				Edge: *graphEdge,
			}
		}
		result.Edges[dataKey] = edges
	}
	return result
}
