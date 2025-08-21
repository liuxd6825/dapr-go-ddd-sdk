package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type TagType struct {
	xbase.BaseModel `bson:",inline"`
	Name            string `gorm:"name" json:"name,omitempty" bson:"name"`
	ParentId        string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`
	IsETag          bool   `gorm:"is_e_tag" json:"isETag,omitempty" bson:"is_e_tag"` //是否企业标签
}
