package query

type MongoQuery struct {
	Collection string           `json:"collection" validate:"required" title:"集合名称"`
	Pipeline   []map[string]any `json:"pipeline" validate:"required" title:"管道"`
}
