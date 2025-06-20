package entity

type Document struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	Text     string `json:"text"`
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
}
