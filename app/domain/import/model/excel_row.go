package model

type ExcelRow struct {
	Id       string         `json:"id" gorm:"id" bson:"id"`
	TenantId string         `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	CaseId   string         `json:"caseId" gorm:"case_id"  bson:"case_id" index:""`
	DocId    string         `json:"docId" gorm:"doc_id" bson:"doc_id" index:""`
	FileId   string         `json:"fileId" gorm:"file_id" bson:"file_id" index:""`
	SheetId  string         `json:"sheetId" gorm:"sheet_id" bson:"sheet_id" index:""`
	RowNum   int64          `json:"rowNum" gorm:"row_num" bson:"row_num" index:""`
	Values   map[string]any `json:"values" gorm:"values;json"   bson:"values"`
}

func (f *ExcelRow) GetTenantId() string {
	return f.TenantId
}

func (f *ExcelRow) SetTenantId(v string) {
	f.TenantId = v
}

func (f *ExcelRow) GetId() string {
	return f.Id
}

func (f *ExcelRow) SetId(v string) {
	f.Id = v
}
