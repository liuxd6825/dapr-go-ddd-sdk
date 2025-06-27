package model

type Relation struct {
	NId     string `json:"nid" bson:"nid"`
	Id      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	CaseId  string `json:"caseId" bson:"case_id"`
	DrawId  string `json:"drawId" bson:"draw_id"`
	RelType string `json:"relType" bson:"rel_type"`
	StartId string `json:"startId" bson:"start_id"`
	EndId   string `json:"endId" bson:"end_id"`
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
