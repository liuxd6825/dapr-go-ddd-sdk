package query

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// RecordFindByIdQuery 按ID查询命令
type RecordFindByIdQuery = ddd_query.FindByIdQuery

// RecordFindByIdsQuery 按多个ID查询命令
type RecordFindByIdsQuery = ddd_query.FindByIdsQuery

// RecordFindAllQuery 查询所有命令
type RecordFindAllQuery struct {
	TenantId string `json:"tenantId"`
}

// RecordFindByCaseIdQuery 按聚合根ID查询命令
type RecordFindByCaseIdQuery struct {
	ddd_query.FindPagingQuery
	CaseId string `json:"caseId" query:"case-id"`
}

// RecordFindByDocIdQuery 按聚合根ID查询命令
type RecordFindByDocIdQuery struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	DocId    string `json:"docId"`
}

type RecordFindByTaskIdQuery struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	TaskId   string `json:"taskId"`
}

// RecordFindByFileIdQuery 按聚合根ID查询命令
type RecordFindByFileIdQuery struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	FileId   string `json:"fileId"`
}

type RecordFindPagingQuery = ddd_query.FindPagingQuery
type RecordFindPagingResult = ddd_query.FindPagingResult

func (q *RecordFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("CaseId=='%s'", q.CaseId), nil
}
