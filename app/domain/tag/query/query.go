package query

type FindPagingTagQuery struct {
	TenantId string `json:"tenantId" query:"tenant-id" required:"true"`
	CaseId   string `json:"caseId"  query:"case-id" required:"-"`
	ETag     bool   `json:"eTag" query:"e-tag" required:"true"`
}

type FindPagingByCaseIdQuery struct {
	CaseId string `json:"caseId"  query:"case-id" required:"true"`
}
