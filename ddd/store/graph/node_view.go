package graph

type Node struct {
	Nid   string            `json:"nid,omitempty"`
	Id    string            `json:"id,omitempty"`
	Label string            `json:"label,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
	Props map[string]any    `json:"props,omitempty"`
}

func (n *Node) GetId() string {
	return n.Id
}

func (n *Node) GetNid() string {
	return n.Nid
}

func (n *Node) GetLabel() string {
	return n.Label
}

func (n *Node) GetTags() map[string]string {
	return n.Tags
}

func (n *Node) GetProps() map[string]any {
	return n.Props
}
