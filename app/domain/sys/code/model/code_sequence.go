package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CodeSequence struct {
	xbase.BaseModel
	TypeCode string `json:"typeCode" bson:"type_code" gorm:"type_code:varchar(255)"`
	DateStr  string `json:"dateStr" bson:"date_str" gorm:"date_str"`
	MaxSeq   int64  `json:"maxSeq" bson:"max_seq" gorm:"max_seq"`
}
