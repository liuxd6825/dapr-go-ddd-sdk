package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

// TaskFindPagingByCaseIdQuery 分页查询命令
type TaskFindPagingByCaseIdQuery struct {
	CaseId string `json:"caseId"`
	idao.FindPagingQueryRequest
}

// TaskFindByCaseIdQuery
// @Description: 根据caseId分页查询
type TaskFindByCaseIdQuery = ddd_query.FindPagingByCaseIdQuery

type SetImportProgressResult struct {
	DataTotal  int64 `json:"dataTotal"`
	ErrorTotal int64 `json:"errorTotal"`
}
