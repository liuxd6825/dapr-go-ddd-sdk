package model

type Dict struct {
	BaseModel    `bson:",inline"`
	CaseTypeId   string `json:"caseTypeId" gorm:"case_type_id" bson:"case_type_id"`
	CaseTypeName string `json:"caseTypeName" gorm:"case_type_name" bson:"case_type_name"`
	CaseName     string `json:"caseName" gorm:"case_name" bson:"case_name"`
	DictTypeId   string `json:"dictTypeId" gorm:"dict_type_id" bson:"dict_type_id"`
	DictTypeName string `json:"dictTypeName" gorm:"dict_type_name" bson:"dict_type_name"`
	DictTypeCode string `json:"dictTypeCode" gorm:"dict_type_code" bson:"dict_type_code"`
	Code         string `json:"code" gorm:"code" bson:"code"`
	Title        string `json:"title" gorm:"title" bson:"title"`
	OrderNum     *int64 `json:"orderNum" gorm:"order_num"  bson:"order_num"`
}

func NewDict() (*Dict, error) {
	return &Dict{}, nil
}
