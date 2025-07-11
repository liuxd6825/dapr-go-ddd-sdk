package model

type File struct {
	Id       string       `json:"id" gorm:"id" bson:"id" `
	TenantId string       `json:"tenantId" gorm:"tenant_id" bson:"tenant_id" index:""`
	CaseId   string       `json:"caseId" gorm:"case_id" bson:"case_id" index:""`
	DocId    string       `json:"docId" gorm:"doc_id" bson:"doc_id" index:""`
	Name     string       `json:"name" gorm:"name" bson:"name" index:""  `
	Sheets   []*SheetInfo `json:"sheets" gorm:"sheets" bson:"sheets"`
}

type SheetInfo struct {
	Name   string `json:"name" bson:"name"`
	MaxCol int64  `json:"maxCol" bson:"max_col"`
	MaxRow int64  `json:"maxRow" bson:"max_row"`
}

func (f *File) GetTenantId() string {
	return f.TenantId
}

func (f *File) SetTenantId(v string) {
	f.TenantId = v
}

func (f *File) GetId() string {
	return f.Id
}

func (f *File) SetId(v string) {
	f.Id = v
}

func (f *File) GetSheet(sheetName string) *SheetInfo {
	for _, s := range f.Sheets {
		if s.Name == sheetName {
			return s
		}
	}
	return nil
}
