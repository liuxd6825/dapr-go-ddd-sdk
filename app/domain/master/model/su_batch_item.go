package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"

type SuBatchItem struct {
	xbase.BaseModel `bson:",inline"`
	BatchId         string `json:"batchId"  gorm:"batch_id"  bson:"batch_id" title:"交易ID"`
	SuRecordId      string `json:"suRecordId" gorm:"su_record_id"  bson:"su_record_id"  index:"" validate:"-" title:"可疑交易ID"`
	RecordId        string `json:"recordId" gorm:"record_id"  bson:"record_id"  index:"" validate:"-" title:"交易ID"`
}
