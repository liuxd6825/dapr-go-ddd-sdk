package entity

type Document struct {
	Id       string `json:"id"`
	Source   string `json:"source"`
	Text     string `json:"text"`
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
}
