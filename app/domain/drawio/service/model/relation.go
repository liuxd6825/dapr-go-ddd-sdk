package model

type Relation struct {
	NId     string `json:"nid"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	CaseId  string `json:"caseId"`
	DrawId  string `json:"drawId"`
	RelType string `json:"relType"`
	StartId string `json:"startId"`
	EndId   string `json:"endId"`
}

func NewRelation() *Relation {
	return &Relation{}
}

func (n *Relation) GetId() string {
	return n.Id
}

type RelationView struct {
	NId     string         `json:"nid"`
	Id      string         `json:"id"`
	RelType string         `json:"relType"`
	StartId string         `json:"startId"`
	EndId   string         `json:"endId"`
	Props   map[string]any `json:"props"`
}

func NewRelationView() *RelationView {
	return &RelationView{}
}
