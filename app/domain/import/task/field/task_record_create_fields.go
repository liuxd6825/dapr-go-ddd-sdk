package field

type TaskRecordCreateFields struct {
	Id        string `json:"id" desc:"Id"`
	BatchSize int64  `json:"batchSize"  desc:"批数据"`
}
