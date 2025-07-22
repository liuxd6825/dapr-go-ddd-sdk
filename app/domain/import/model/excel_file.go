package model

type ExcelFile struct {
	Id        string `json:"id" gorm:"id" bson:"id" `
	Name      string `json:"name"  gorm:"name" bson:"name"  index:""  `
	TenantId  string `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	CaseId    string `json:"caseId"  gorm:"case_id" bson:"case_id" index:""`
	DocId     string `json:"docId"  gorm:"doc_id" bson:"doc_id" index:""`
	DocFileId string `json:"docFileId" gorm:"doc_file_id" bson:"doc_file_id"`
}

func (f *ExcelFile) GetTenantId() string {
	return f.TenantId
}

func (f *ExcelFile) SetTenantId(v string) {
	f.TenantId = v
}

func (f *ExcelFile) GetId() string {
	return f.Id
}

func (f *ExcelFile) SetId(v string) {
	f.Id = v
}
