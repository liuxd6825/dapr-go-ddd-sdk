package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type Tag struct {
	xbase.BaseModel `bson:",inline"`
	TagTypeId       string `gorm:"tag_type_id" json:"tagTypeId,omitempty" bson:"tag_type_id"`
	Name            string `gorm:"name" json:"name,omitempty" bson:"name"`
	Color           string `gorm:"color" json:"color,omitempty" bson:"color"`
	IsETag          bool   `gorm:"is_e_tag" json:"isETag,omitempty" bson:"is_e_tag"` //是否企业标签
}
