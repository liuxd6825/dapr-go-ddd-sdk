package core

type Node struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	CaseId      string   `json:"caseId"`
	TenantId    string   `json:"tenantId"`
	Keywords    []string `json:"keywords"`
	Table       string   `json:"table"`
	Description string   `json:"description"`
	SourceIds   []string `json:"source_ids"`
	SourceType  string   `json:"source_type"`
	SourceName  string   `json:"source_name"`
	SourceUrl   string   `json:"source_url"`
}

/*
	CREATE (n$<nLabels>{$<props>}) WITH n
	MATCH (m$<mLabels>{id:$<source>}) WITH m, n
	CREATE (m)-[r:$<relType>]->(n)
	SET r.id=$<id>,r.keywords=$<keywords>,r.case_id=$<caseId>,r.tenant_id=$<tenantId>,
		r.description=$<description>,r.source=$<source>,r.target=$<target>,
		r.source_ids=$<id>,r.source_type=$<sourceType>,r.table=$<table>
*/
