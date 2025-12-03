package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CodeType struct {
	xbase.BaseModel
	Code       string `json:"code" bson:"code" gorm:"code"`                     // e.g., "MH"
	Prefix     string `json:"prefix" bson:"prefix" gorm:"prefix"`               // e.g., "MH-"
	DateFormat string `json:"dateFormat" bson:"date_format" gorm:"date_format"` // e.g., "060102" (YYMMDD)
	SeqLength  int    `json:"seqLength" bson:"seq_length" json:"seq_length"`    // e.g., 6
}

type CodeStyle string

const (
	CodeStyle_Global CodeStyle = ""
	CodeStyle_Year   CodeStyle = "year"
	CodeStyle_Month  CodeStyle = "month"
	CodeStyle_Day    CodeStyle = "day"
)

func (db CodeStyle) String() string {
	return string(db)
}

func (db CodeStyle) DateFormat() string {
	switch db {
	case CodeStyle_Year:
		return "06"
	case CodeStyle_Month:
		return "0601"
	case CodeStyle_Day:
		return "060102"
	default:
		return ""
	}
}
