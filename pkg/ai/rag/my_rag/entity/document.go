package entity

type Document struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	Text     string `json:"text"`
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	Url      string `json:"url"`
}

func (d *Document) GetId() string {
	return d.Id
}

func (d *Document) GetCaseId() string {
	return d.CaseId
}

func (d *Document) GetText() string {
	return d.Text
}

func (d *Document) GetFileName() string {
	return d.FileName
}

func (d *Document) GetTenantId() string {
	return d.TenantId
}

func (d *Document) GetUrl() string {
	return d.Url
}
