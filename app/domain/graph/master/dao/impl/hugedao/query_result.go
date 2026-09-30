package hugedao

type QueryResult struct {
	Data []QueryResultItem
}

type QueryResultItem struct {
	Id         string            `json:"id"`
	Label      string            `json:"label"`
	Properties map[string]string `json:"properties"`
}
