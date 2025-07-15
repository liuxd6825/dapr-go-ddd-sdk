package store

import (
	"errors"
	"fmt"
)

type FindPagingByCaseIdQuery interface {
	FindPagingQuery
	GetCaseId() string
	SetCaseId(v string)
}

type FindPagingByCaseIdQueryRequest struct {
	FindPagingQueryRequest
	CaseId string `json:"caseId" required:"true" query:"case-id" path:"caseId" title:"案件ID" `
}

func NewFindPagingByCaseIdQuery(paging *FindPagingQueryRequest, caseId string) *FindPagingByCaseIdQueryRequest {
	return &FindPagingByCaseIdQueryRequest{
		FindPagingQueryRequest: *paging,
		CaseId:                 caseId,
	}
}

func (q *FindPagingByCaseIdQueryRequest) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("CaseId不能为空")
	}
	return fmt.Sprintf("caseId=='%s'", q.CaseId), nil
}

func (q *FindPagingByCaseIdQueryRequest) GetCaseId() string {
	return q.CaseId
}

func (q *FindPagingByCaseIdQueryRequest) SetCaseId(v string) {
	q.CaseId = v
}
