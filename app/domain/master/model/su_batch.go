package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type SuBatch struct {
	xbase.BaseModel `bson:",inline"`
	Name            string  `json:"name" gorm:"name" bson:"name" title:"名称"`
	Type            SuType  `json:"type" gorm:"type" bson:"type" title:"类型"`
	Count           int     `json:"count" gorm:"count" bson:"count" title:"交易数量"`
	Payout          float64 `json:"payout" gorm:"payout" bson:"payout" title:"支出总额"`
	Income          float64 `json:"income" gorm:"income" bson:"income" title:"收入总额"`
}
