package model

type Node struct {
	Nid      string `json:"nid" bson:"nid"`
	Id       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Label    string `json:"label" bson:"label"`
	CaseId   string `json:"caseId" bson:"case_id"`
	TenantId string `json:"tenantId" bson:"tenant_id"`
	DrawId   string `json:"drawId" bson:"draw_id"`
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
