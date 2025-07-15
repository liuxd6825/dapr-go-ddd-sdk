package model

type ExcelRow struct {
	Id       string         `json:"id" bson:"id"`
	TenantId string         `json:"tenantId" bson:"tenant_id" index:""`
	CaseId   string         `json:"caseId" bson:"case_id" index:""`
	DocId    string         `json:"docId" bson:"doc_id" index:""`
	FileId   string         `json:"fileId" bson:"file_id" index:""`
	SheetId  string         `json:"sheetId" bson:"sheet_id" index:""`
	RowNum   int64          `json:"rowNum" bson:"row_num" index:""`
	Values   map[string]any `json:"values" bson:"values"`
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
