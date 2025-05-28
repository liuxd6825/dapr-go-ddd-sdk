package model

type Relation struct {
	Id      string `json:"id"`
	CaseId  string `json:"caseId"`
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
