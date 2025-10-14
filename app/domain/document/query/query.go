package query

type FindByIdQuery struct {
	Id string `json:"id" param:"id" required:"true"`
}

type FindFolderByFolderIdQuery struct {
	FolderId string `json:"folderId" query:"folder-id" required:"true"`
}

type FindByCaseIdQuery struct {
	CaseId string `json:"caseId"  query:"case-id" required:"true"`
}

type FindByDocumentIdQuery struct {
	DocumentId string `json:"documentId"  query:"document-id" required:"true"`
}

type FindTreeByParamQuery struct {
	TenantId string `json:"tenantId" query:"tenant-id" required:"true"`
	BusId    string `json:"busId" query:"bus-id" required:"true"`
	EntityId string `json:"entityId" query:"entity-id" required:"true"`
}

type DownloadParamQuery struct {
	FileId string `json:"fileId" query:"file-id" required:"true"`
}

type FindByFolderAndFilter struct {
	FolderId string `json:"folderId" query:"folder-id" required:"true"`
	Filter   string `json:"filter" query:"filter" required:"-"`
}

type FindByTagAndCase struct {
	CaseId string `json:"caseId"  query:"case-id" required:"true"`
	TagId  string `json:"tagId"  query:"tag-id" required:"true"`
}
