package model

type Relation struct {
	NId         string   `json:"nid" bson:"nid"`
	Id          string   `json:"id" bson:"id" title:"ID"`
	Name        string   `json:"name" bson:"name" title:"名称"`
	CaseId      string   `json:"caseId" bson:"case_id" title:"案件ID"`
	DrawId      string   `json:"drawId" bson:"draw_id" title:"绘图ID"`
	RelType     string   `json:"relType" bson:"rel_type" title:"关系类型"`
	Source      string   `json:"source" bson:"source" title:"来源ID"`
	Target      string   `json:"target" bson:"target" title:"目录ID"`
	Keywords    []string `json:"keywords" bson:"keywords" title:"关键字"`
	SourceIds   string   `json:"sourceIds" bson:"source_ids" title:"数据源ID"`
	SourceType  string   `json:"sourceType" bson:"source_type" title:"数据源类型"`
	SourceUrl   string   `json:"sourceUrl" gorm:"source_url"  bson:"source_url" title:"数据源URL"`
	SourceName  string   `json:"sourceName" gorm:"source_name"  bson:"source_name" title:"数据源名称"`
	Description string   `json:"description" bson:"description" title:"说明"`
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
