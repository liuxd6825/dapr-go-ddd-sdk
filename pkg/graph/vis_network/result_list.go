package vis_network

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
)

type NodeList []*Node
type EdgeList []*Edge

// ResultList Vis数据
type ResultList struct {
	Nodes     NodeList       `json:"nodes"`
	Edges     EdgeList       `json:"edges"`
	Header    map[string]any `json:"header"`
	PageTotal *int64         `json:"pageTotal"`
	PageSize  *int64         `json:"pageSize"`
	PageNum   *int64         `json:"pageNum"`
	Total     *int64         `json:"total"`
	nodeMap   NodeMap
	edgeMap   EdgeMap
}

func NewResultList() *ResultList {
	return &ResultList{
		Header:  nil,
		Nodes:   NodeList{},
		Edges:   EdgeList{},
		nodeMap: NodeMap{},
		edgeMap: EdgeMap{},
	}
}

type ResultListWithGraphViewOption struct {
	OnNode func(result *ResultList, n *Node)
	OnEdge func(result *ResultList, e *Edge)
}

func NewResultListWithGraphViewOption(opts ...*ResultListWithGraphViewOption) *ResultListWithGraphViewOption {
	opt := &ResultListWithGraphViewOption{}
	for _, o := range opts {
		if o != nil {
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

func NewResultListWithGraphView(graphView *graph.GraphView, opts ...*ResultListWithGraphViewOption) *ResultList {
	r := NewResultList()
	opt := NewResultListWithGraphViewOption(opts...)

	for _, graphNodes := range graphView.Nodes {
		r.Nodes = r.addNodeList(graphNodes, r.Nodes, opt)
	}
	for _, graphEdges := range graphView.Edges {
		r.Edges = r.addEdgeList(graphEdges, r.Edges, opt)
	}
	return r
}

func (r *ResultList) addNodeList(graphNodes graph.Nodes, nodeList NodeList, opt *ResultListWithGraphViewOption) NodeList {
	for key, graphNode := range graphNodes {
		if _, ok := r.nodeMap[key]; ok {
			continue
		}
		n := &Node{
			Node: *graphNode,
		}
		r.nodeMap[key] = n
		nodeList = append(nodeList, n)
		if opt.OnNode != nil {
			opt.OnNode(r, n)
		}
	}
	return nodeList
}

func (r *ResultList) addEdgeList(graphEdges graph.Edges, edgeList EdgeList, opt *ResultListWithGraphViewOption) EdgeList {
	for key, graphEdge := range graphEdges {
		if _, ok := r.edgeMap[key]; ok {
			continue
		}
		e := &Edge{
			Edge: *graphEdge,
		}
		edgeList = append(edgeList, e)
		r.edgeMap[key] = e
		if opt.OnEdge != nil {
			opt.OnEdge(r, e)
		}
	}
	return edgeList
}
