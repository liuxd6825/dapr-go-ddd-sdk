package graph

type Nodes map[string]*Node
type Edges map[string]*Edge

type GraphView struct {
	Nodes map[string]Nodes
	Edges map[string]Edges
}

func NewGraphView() *GraphView {
	return &GraphView{
		Nodes: make(map[string]Nodes),
		Edges: make(map[string]Edges),
	}
}

func (g *GraphView) SetNodes(key string, data map[string]*Node) {
	g.Nodes[key] = data
}

func (g *GraphView) SetEdges(key string, data map[string]*Edge) {
	g.Edges[key] = data
}

func (g *GraphView) AddNodes(key string, data []*Node) {
	items, ok := g.Nodes[key]
	if !ok {
		items = make(map[string]*Node)
	}
	for _, node := range data {
		items[node.Id] = node
	}
	g.Nodes[key] = items
}

func (g *GraphView) AddEdges(key string, data []*Edge) {
	items, ok := g.Edges[key]
	if !ok {
		items = make(map[string]*Edge)
	}
	for _, item := range data {
		items[item.Id] = item
	}
	g.Edges[key] = items
}
