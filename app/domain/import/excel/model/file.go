package model

type ExcelFile struct {
	Id       string                `json:"id" gorm:"id" bson:"id" `
	TenantId string                `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	CaseId   string                `json:"caseId"  gorm:"case_id" bson:"case_id" index:""`
	DocId    string                `json:"docId"  gorm:"doc_id" bson:"doc_id" index:""`
	Name     string                `json:"name"  gorm:"name" bson:"name"  index:""  `
	Sheets   []*ExcelFileSheetInfo `json:"sheets"  gorm:"sheets" bson:"sheets"`
}

type ExcelFileSheetInfo struct {
	Name   string `json:"name"`
	MaxCol int64  `json:"maxCol"`
	MaxRow int64  `json:"maxRow"`
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

func (f *ExcelFile) GetSheet(sheetName string) *ExcelFileSheetInfo {
	for _, s := range f.Sheets {
		if s.Name == sheetName {
			return s
		}
	}
	return nil
}
