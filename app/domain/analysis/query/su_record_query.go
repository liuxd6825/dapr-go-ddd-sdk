package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"

// SuRecordFindByIdQuery 按聚合根ID查询命令
type SuRecordFindByIdQuery struct {
	Id string `json:"id" param:"id"`
}

type SuRecordFindByCaseIdQuery struct {
	ddd_query.FindPagingQuery
	CaseId string `json:"caseId" query:"case-id"`
	TaskId string `json:"taskId" query:"task-id"`
}
