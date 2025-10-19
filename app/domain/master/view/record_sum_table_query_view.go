package view

type RecordSumTableQueryView struct {
	Data      []*RecordSumTableView `json:"data" bson:"data"`
	TotalRows int64                 `json:"totalRows" bson:"total_rows"`
}
