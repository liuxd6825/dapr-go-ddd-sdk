package model

type Node struct {
	Nid         string   `json:"nid" gorm:"nid"  bson:"nid"`
	Id          string   `json:"id"  gorm:"id" bson:"id"`
	Name        string   `json:"name" gorm:"name"  bson:"name"`
	Type        string   `json:"type" gorm:"type"  bson:"type"`
	CaseId      string   `json:"caseId" gorm:"case_id"  bson:"case_id"`
	TenantId    string   `json:"tenantId"  gorm:"tenant_id" bson:"tenant_id"`
	DrawId      string   `json:"drawId" gorm:"draw_id"  bson:"draw_id"`
	SrcId       string   `json:"srcId" gorm:"src_id"  bson:"src_id"`
	SrcType     string   `json:"srcType" gorm:"src_type"  bson:"src_type"`
	SrcUrl      string   `json:"srcUrl" gorm:"src_url"  bson:"src_url" title:"数据源URL"`
	SrcName     string   `json:"srcName" gorm:"src_name"  bson:"src_name" title:"数据源名称"`
	Description string   `json:"description" gorm:"description" bson:"description"`
	Keywords    []string `json:"keywords" gorm:"keywords"  bson:"keywords"`
	Table       string   `json:"table" gorm:"table"  bson:"table"`
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
