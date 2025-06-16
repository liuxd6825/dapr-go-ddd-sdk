package model

type Message struct {
	Base     `bson:",inline"`
	Content  string `gorm:"title" json:"title,omitempty" bson:"content"` // 内容
	ParentId string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`
}
