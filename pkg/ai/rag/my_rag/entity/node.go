package entity

type Node1 struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	CaseId      string `json:"case_id"`
	TenantId    string `json:"tenant_id"`
	DocId       string `json:"doc_id"`
	Description string `json:"description"`
	SourceIds   string `json:"source_ids"`
	SourceType  string `json:"source_type"`
	SourceUrl   string `json:"source_url"`
	SourceName  string `json:"source_name"`
}
