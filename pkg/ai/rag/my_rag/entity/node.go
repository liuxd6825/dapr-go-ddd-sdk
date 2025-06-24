package entity

type Node struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CaseId      string `json:"case_id"`
	TenantId    string `json:"tenant_id"`
	DocId       string `json:"doc_id"`
	Description string `json:"description"`
	SourceIds   string `json:"source_ids"`
	SourceType  string `json:"source_type"`
	Type        string `json:"type"`
}
