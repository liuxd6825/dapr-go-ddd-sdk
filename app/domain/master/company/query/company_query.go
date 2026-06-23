package query

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// CompanyFindByIdQuery 按ID查询命令
type CompanyFindByIdQuery = store.FindByIdRequest

// CompanyFindPagingQuery 分页查询命令
type CompanyFindPagingQuery = store.FindPagingQueryRequest

// CompanyFindByCaseIdQuery 按案件ID分页查询命令
type CompanyFindByCaseIdQuery struct {
	store.FindPagingQueryRequest
	CaseId string `json:"caseId" query:"case-id"`
}

func (q *CompanyFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("case_id=='%s'", q.CaseId), nil
}

// CompanyFindByTagIdQuery 按案件ID和标签ID查询命令
type CompanyFindByTagIdQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	TagId  string `json:"tagId" query:"tag-id" validate:"required" title:"标签ID"`
}

// CompanyFindTotalQuery 统计查询命令
type CompanyFindTotalQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	Icon   string `json:"icon" query:"icon" title:"图标"`
}

// CompanyFindTotalResult 统计结果
type CompanyFindTotalResult struct {
	Count int64  `json:"count" title:"公司总数"`
	Icon  string `json:"icon" title:"图标"`
}