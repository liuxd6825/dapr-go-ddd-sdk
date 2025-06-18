package model

type Message struct {
	Base     `bson:",inline"`
	ChatId   string  `gorm:"chat_id" json:"chatId" bson:"chat_id"`
	Content  string  `gorm:"content;size:8000" json:"content,omitempty" bson:"content"` // 内容
	ParentId string  `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`
	OrderNum float32 `gorm:"order_num" json:"orderNum,omitempty" bson:"order_num"`
	Role     string  `gorm:"role" json:"role,omitempty" bson:"role"`
}
