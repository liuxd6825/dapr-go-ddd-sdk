package query

// SuTaskFindByIdQuery 按聚合根ID查询命令
type SuTaskFindByIdQuery struct {
	Id string `json:"id" param:"id"`
}

type SuTaskFindByCaseIdQuery struct {
	CaseId string `json:"caseId" query:"case-id"`
}
