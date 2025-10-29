package model

type DictType struct {
	BaseModel `bson:",inline"`
	Name      string `json:"name"  gorm:"name" bson:"name"`
	Code      string `json:"code" gorm:"code"  bson:"code"`
	Category  string `json:"category" gorm:"category"  bson:"category"`
	ParentId  string `json:"parentId" gorm:"parent_id" bson:"parent_id"`
}

func NewDictType() (*DictType, error) {
	return &DictType{}, nil
}
