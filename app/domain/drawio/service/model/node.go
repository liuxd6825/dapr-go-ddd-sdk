package model

type Node struct {
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
