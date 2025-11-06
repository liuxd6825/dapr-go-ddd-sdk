package graph

type Edge struct {
	Nid   string         `json:"nid,omitempty"`
	NFrom string         `json:"nfrom,omitempty"`
	NTo   string         `json:"nto,omitempty"`
	Id    string         `json:"id"`
	Label string         `json:"label"`
	From  string         `json:"from,omitempty"`
	To    string         `json:"to,omitempty"`
	Props map[string]any `json:"props"`
}

func (n *Edge) GetId() string {
	return n.Id
}

func (n *Edge) GetNid() string {
	return n.Nid
}

func (n *Edge) GetName() string {
	return n.Label
}

func (n *Edge) GetFrom() string {
	return n.From
}

func (n *Edge) GetTo() string {
	return n.To
}

func (n *Edge) GetProps() map[string]any {
	return n.Props
}
