package model

type ExcelSheet struct {
	Id       string `json:"id" bson:"id"`
	TenantId string `json:"tenantId" bson:"tenant_id" index:""`
	Name     string `json:"name" bson:"name" index:""`

	CaseId   string `json:"caseId" bson:"case_id" index:""`
	DocId    string `json:"docId" bson:"doc_id" index:""`
	FileId   string `json:"fileId" bson:"file_id"  index:""`
	FileName string `json:"fileName" bson:"file_name" index:""`

	MaxRow  int64    `json:"maxRow" bson:"max_row"`
	MaxCol  int64    `json:"maxCol" bson:"max_col"`
	Columns []string `json:"columns" bson:"columns"`
}

func (f *ExcelSheet) GetTenantId() string {
	return f.TenantId
}

func (f *ExcelSheet) SetTenantId(v string) {
	f.TenantId = v
}

func (f *ExcelSheet) GetId() string {
	return f.Id
}

func (f *ExcelSheet) SetId(v string) {
	f.Id = v
}
