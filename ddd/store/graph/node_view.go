package graph

type Node struct {
	Nid   string         `json:"nid,omitempty"`
	Id    string         `json:"id,omitempty"`
	Name  string         `json:"name,omitempty"`
	Tags  []string       `json:"tags,omitempty"`
	Props map[string]any `json:"props,omitempty"`
}

func (n *Node) GetId() string {
	return n.Id
}

func (n *Node) GetNid() string {
	return n.Nid
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) GetTags() []string {
	return n.Tags
}

func (n *Node) GetProps() map[string]any {
	return n.Props
}
