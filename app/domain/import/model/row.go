package model

type Row struct {
	Id       string         `json:"id" gorm:"id" bson:"id"`
	TenantId string         `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	CaseId   string         `json:"caseId" gorm:"case_id" bson:"case_id" index:""`
	DocId    string         `json:"docId" gorm:"doc_id" bson:"doc_id" index:""`
	FileId   string         `json:"fileId" gorm:"file_id" bson:"file_id" index:""`
	SheetId  string         `json:"sheetId" gorm:"sheet_id" bson:"sheet_id" index:""`
	RowNum   int64          `json:"rowNum" gorm:"row_num" bson:"row_num" index:""`
	Values   map[string]any `json:"values" gorm:"values" bson:"values"`
}

func (f *Row) GetTenantId() string {
	return f.TenantId
}

func (f *Row) SetTenantId(v string) {
	f.TenantId = v
}

func (f *Row) GetId() string {
	return f.Id
}

func (f *Row) SetId(v string) {
	f.Id = v
}
