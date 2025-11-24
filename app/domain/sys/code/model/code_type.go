package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"

type CodeType struct {
	xbase.BaseModel
	Code       string `json:"code" bson:"code" gorm:"code"`                     // e.g., "MH"
	Prefix     string `json:"prefix" bson:"prefix" gorm:"prefix"`               // e.g., "MH-"
	DateFormat string `json:"dateFormat" bson:"date_format" gorm:"date_format"` // e.g., "060102" (YYMMDD)
	SeqLength  int    `json:"seqLength" bson:"seq_length" json:"seq_length"`    // e.g., 6
}
