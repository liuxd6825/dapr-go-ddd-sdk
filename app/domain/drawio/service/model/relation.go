package model

type Relation struct {
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
