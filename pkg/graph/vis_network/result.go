package vis_network

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"

type NodeMap map[string]*Node
type EdgeMap map[string]*Edge

// Result Vis数据
type Result struct {
	Nodes     map[string]NodeMap `json:"nodes"`
	Edges     map[string]EdgeMap `json:"edges"`
	Header    map[string]any     `json:"header"`
	PageTotal *int64             `json:"pageTotal"`
	PageSize  *int64             `json:"pageSize"`
	PageNum   *int64             `json:"pageNum"`
	Total     *int64             `json:"total"`
}

func NewResult() *Result {
	return &Result{
		Header: nil,
		Nodes:  map[string]NodeMap{},
		Edges:  map[string]EdgeMap{},
	}
}

type ResultWithGraphViewOption struct {
	Merge  *bool //是否对数据进行合并去重处理
	OnNode func(result *Result, n *Node)
	OnEdge func(result *Result, e *Edge)
}

func NewResultWithGraphViewOption(opts ...*ResultWithGraphViewOption) *ResultWithGraphViewOption {
	opt := &ResultWithGraphViewOption{}
	for _, o := range opts {
		if o != nil {
			if o.Merge != nil {
				opt.Merge = o.Merge
			}
			if o.OnNode != nil {
				opt.OnNode = o.OnNode
			}
			if o.OnEdge != nil {
				opt.OnEdge = o.OnEdge
			}
		}
	}
	return opt
}

func NewResultWithGraphView(graphView *graph.GraphView, opts ...*ResultWithGraphViewOption) *Result {
	result := &Result{
		Header: nil,
		Nodes:  map[string]NodeMap{},
		Edges:  map[string]EdgeMap{},
	}
	opt := NewResultWithGraphViewOption(opts...)
	if opt.Merge != nil && *opt.Merge {
		nodeMap := NodeMap{}
		result.Nodes["n"] = nodeMap
		for _, graphNodes := range graphView.Nodes {
			addNodeMap(result, graphNodes, nodeMap, opt)
		}
		edgeMap := EdgeMap{}
		result.Edges["r"] = edgeMap
		for _, graphEdges := range graphView.Edges {
			addEdgeMap(result, graphEdges, edgeMap, opt)
		}
		return result
	}

	for dataKey, graphNodes := range graphView.Nodes {
		nodeMap := NodeMap{}
		result.Nodes[dataKey] = nodeMap
		addNodeMap(result, graphNodes, nodeMap, opt)
	}
	for dataKey, graphEdges := range graphView.Edges {
		edgeMap := EdgeMap{}
		result.Edges[dataKey] = edgeMap
		addEdgeMap(result, graphEdges, edgeMap, opt)
	}
	return result
}

func addNodeMap(result *Result, graphNodes graph.Nodes, nodeMap NodeMap, opt *ResultWithGraphViewOption) {
	for key, graphNode := range graphNodes {
		n := &Node{
			Node: *graphNode,
		}
		nodeMap[key] = n
		if opt.OnNode != nil {
			opt.OnNode(result, n)
		}
	}
}

func addEdgeMap(result *Result, graphEdges graph.Edges, edgeMap EdgeMap, opt *ResultWithGraphViewOption) {
	for key, graphEdge := range graphEdges {
		e := &Edge{
			Edge: *graphEdge,
		}
		edgeMap[key] = e
		if opt.OnEdge != nil {
			opt.OnEdge(result, e)
		}
	}
}
