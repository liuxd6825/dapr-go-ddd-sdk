package entity

type Edge struct {
	Id       string `json:"id"`
	CaseId   string `json:"case_id"`
	TenantId string `json:"tenant_id"`
	DocId    string `json:"doc_id"`

	Target      string `json:"target"`
	Source      string `json:"source"`
	Description string `json:"description"`
	SourceIds   string `json:"source_ids"`
	SourceType  string `json:"source_type"`
	Keywords    string `json:"keywords"`
	RelType     string `json:"rel_type"`
}
