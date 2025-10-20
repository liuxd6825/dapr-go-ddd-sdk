package model

type Login struct {
	Method     string `json:"method" gorm:"method" bson:"method"`
	Identifier string `json:"identifier" gorm:"identifier" bson:"identifier"`
	Password   string `json:"password" gorm:"password" bson:"password"`
	Flow       string `json:"flow" gorm:"flow" bson:"flow"`
}
