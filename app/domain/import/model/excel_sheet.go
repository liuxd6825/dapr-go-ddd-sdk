package model

type ExcelSheet struct {
	Id       string `json:"id" gorm:"id" bson:"id"`
	TenantId string `json:"tenantId" gorm:"tenant_id"  bson:"tenant_id" index:""`
	Name     string `json:"name" gorm:"name"  bson:"name" index:""`

	TaskId    string `json:"taskId" gorm:"task_id"  bson:"task_id" index:""`
	CaseId    string `json:"caseId" gorm:"case_id"  bson:"case_id" index:""`
	DocId     string `json:"docId" gorm:"doc_id"  bson:"doc_id" index:""`
	DocFileId string `json:"docFileId" gorm:"doc_file_id" bson:"doc_file_id"`
	FileId    string `json:"fileId" gorm:"file_id"  bson:"file_id"  index:""`
	FileName  string `json:"fileName" gorm:"file_name"  bson:"file_name" index:""`

	MaxRow  int64    `json:"maxRow" gorm:"max_row"  bson:"max_row"`
	MaxCol  int64    `json:"maxCol" gorm:"max_col"  bson:"max_col"`
	Columns []string `json:"columns" gorm:"columns;json"  bson:"columns"`
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
