package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type TagRelation struct {
	xbase.BaseModel `bson:",inline"`
	TagId           string `gorm:"tag_id" json:"tagId,omitempty" bson:"tag_id"`                         //标签Id
	TagName         string `gorm:"tag_name" json:"tagName,omitempty" bson:"tag_name"`                   //标签名
	TagColor        string `gorm:"tag_color" json:"tagColor,omitempty" bson:"tag_color"`                //标签颜色
	BusType         string `gorm:"bus_type" json:"busType,omitempty" bson:"bus_type"`                   //业务类型
	BusId           string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`                         //业务Id
	ChangedSource   string `gorm:"changed_source" json:"changedSource,omitempty" bson:"changed_source"` //变更来源  TagCenter   BusinessSystem
}
