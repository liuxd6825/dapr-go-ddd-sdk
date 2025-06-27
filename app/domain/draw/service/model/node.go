package model

type Node struct {
	Nid         string `json:"nid" bson:"nid"`
	Id          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	CaseId      string `json:"caseId" bson:"case_id"`
	TenantId    string `json:"tenantId" bson:"tenant_id"`
	DrawId      string `json:"drawId" bson:"draw_id"`
	SourceIds   string `json:"sourceIds" bson:"source_ids"`
	SourceType  string `json:"sourceType" bson:"source_type"`
	Description string `json:"description" bson:"description"`
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
