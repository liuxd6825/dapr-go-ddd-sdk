package model

type Chat struct {
	Base  `bson:",inline"`
	Title string `gorm:"title" json:"title,omitempty" bson:"title"` // 标题
}
