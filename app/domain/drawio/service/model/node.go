package model

type Node struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Labels   []string `json:"labels"`
	CaseId   string   `json:"caseId"`
	TenantId string   `json:"tenantId"`
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) GetId() string {
	return n.Id
}
