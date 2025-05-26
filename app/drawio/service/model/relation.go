package model

type Relation struct {
	Id      string `json:"id"`
	RelType string `json:"relType"`
	StartId string `json:"startId"`
	EndId   string `json:"endId"`
}
