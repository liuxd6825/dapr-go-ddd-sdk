package model

type Relation struct {
	NId         string   `json:"nid" bson:"nid"`
	Id          string   `json:"id" bson:"id"`
	Name        string   `json:"name" bson:"name"`
	CaseId      string   `json:"caseId" bson:"case_id"`
	DrawId      string   `json:"drawId" bson:"draw_id"`
	RelType     string   `json:"relType" bson:"rel_type"`
	Source      string   `json:"source" bson:"source"`
	Description string   `json:"description" bson:"description"`
	Target      string   `json:"target" bson:"target"`
	Keywords    []string `json:"keywords" bson:"keywords"`
	SourceIds   string   `json:"sourceIds" bson:"source_ids"`
	SourceType  string   `json:"sourceType" bson:"source_type"`
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
