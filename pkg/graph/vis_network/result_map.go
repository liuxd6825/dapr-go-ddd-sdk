package vis_network

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"

type NodeMap map[string]*Node
type EdgeMap map[string]*Edge

// ResultMap Vis数据
type ResultMap struct {
	Nodes     map[string]NodeMap `json:"nodes"`
	Edges     map[string]EdgeMap `json:"edges"`
	Header    map[string]any     `json:"header"`
	PageTotal *int64             `json:"pageTotal"`
	PageSize  *int64             `json:"pageSize"`
	PageNum   *int64             `json:"pageNum"`
	Total     *int64             `json:"total"`
}

func NewResultMap() *ResultMap {
	return &ResultMap{
		Header: nil,
		Nodes:  map[string]NodeMap{},
		Edges:  map[string]EdgeMap{},
	}
}

type ResultMapWithGraphViewOption struct {
	Merge  *bool //是否对数据进行合并去重处理
	OnNode func(result *ResultMap, n *Node)
	OnEdge func(result *ResultMap, e *Edge)
}

func NewResultWithGraphViewOption(opts ...*ResultMapWithGraphViewOption) *ResultMapWithGraphViewOption {
	opt := &ResultMapWithGraphViewOption{}
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

func NewResultMapWithGraphView(graphView *graph.GraphView, opts ...*ResultMapWithGraphViewOption) *ResultMap {
	r := &ResultMap{
		Header: nil,
		Nodes:  map[string]NodeMap{},
		Edges:  map[string]EdgeMap{},
	}
	opt := NewResultWithGraphViewOption(opts...)
	if opt.Merge != nil && *opt.Merge {
		nodeMap := NodeMap{}
		r.Nodes["n"] = nodeMap
		for _, graphNodes := range graphView.Nodes {
			r.addNodeMap(graphNodes, nodeMap, opt)
		}
		edgeMap := EdgeMap{}
		r.Edges["r"] = edgeMap
		for _, graphEdges := range graphView.Edges {
			r.addEdgeMap(graphEdges, edgeMap, opt)
		}
		return r
	}

	for dataKey, graphNodes := range graphView.Nodes {
		nodeMap := NodeMap{}
		r.Nodes[dataKey] = nodeMap
		r.addNodeMap(graphNodes, nodeMap, opt)
	}
	for dataKey, graphEdges := range graphView.Edges {
		edgeMap := EdgeMap{}
		r.Edges[dataKey] = edgeMap
		r.addEdgeMap(graphEdges, edgeMap, opt)
	}
	return r
}

func (r *ResultMap) addNodeMap(graphNodes graph.Nodes, nodeMap NodeMap, opt *ResultMapWithGraphViewOption) {
	for key, graphNode := range graphNodes {
		n := &Node{
			Node: *graphNode,
		}
		nodeMap[key] = n
		if opt.OnNode != nil {
			opt.OnNode(r, n)
		}
	}
}

func (r *ResultMap) addEdgeMap(graphEdges graph.Edges, edgeMap EdgeMap, opt *ResultMapWithGraphViewOption) {
	for key, graphEdge := range graphEdges {
		e := &Edge{
			Edge: *graphEdge,
		}
		edgeMap[key] = e
		if opt.OnEdge != nil {
			opt.OnEdge(r, e)
		}
	}
}
