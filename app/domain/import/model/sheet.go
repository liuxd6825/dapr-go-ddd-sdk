package model

type Sheet struct {
	Id       string `json:"id" gorm:"id" bson:"id"`
	TenantId string `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	Name     string `json:"name" gorm:"name" bson:"name" index:""`

	CaseId   string `json:"caseId" gorm:"case_id" bson:"case_id" index:""`
	DocId    string `json:"docId" gorm:"doc_id" bson:"doc_id" index:""`
	FileId   string `json:"fileId" gorm:"file_id" bson:"file_id"  index:""`
	FileName string `json:"fileName" gorm:"file_name" bson:"file_name" index:""`

	MaxRow  int64    `json:"maxRow" gorm:"max_row" bson:"max_row"`
	MaxCol  int64    `json:"maxCol" gorm:"max_col" bson:"max_col"`
	Columns []string `json:"columns" gorm:"columns" bson:"columns"`
}

func (f *Sheet) GetTenantId() string {
	return f.TenantId
}

func (f *Sheet) SetTenantId(v string) {
	f.TenantId = v
}

func (f *Sheet) GetId() string {
	return f.Id
}

func (f *Sheet) SetId(v string) {
	f.Id = v
}
