package model

type Node struct {
	Nid      string `json:"nid"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	CaseId   string `json:"caseId"`
	TenantId string `json:"tenantId"`
	DrawId   string `json:"drawId"`
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) GetId() string {
	return n.Id
}

type NodeView struct {
	Nid   string         `json:"nid"`
	Id    string         `json:"id"`
	Name  string         `json:"name"`
	Tags  []string       `json:"tags"`
	Props map[string]any `json:"props"`
}

func NewNodeView() *NodeView {
	return &NodeView{}
}
