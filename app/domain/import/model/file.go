package model

type File struct {
	Id       string       `json:"id" gorm:"_id" `
	TenantId string       `json:"tenantId" gorm:"tenant_id" index:""`
	CaseId   string       `json:"caseId" gorm:"case_id" index:""`
	DocId    string       `json:"docId" gorm:"doc_id" index:""`
	Name     string       `json:"name" gorm:"name"  index:""  `
	Sheets   []*SheetInfo `json:"sheets" gorm:"sheets"`
}

type SheetInfo struct {
	Name   string `json:"name"`
	MaxCol int64  `json:"maxCol"`
	MaxRow int64  `json:"maxRow"`
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
