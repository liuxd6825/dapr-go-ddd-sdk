package query

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// ProductFindByIdQuery 按ID查询命令
type ProductFindByIdQuery = store.FindByIdRequest

// ProductFindPagingQuery 分页查询命令
type ProductFindPagingQuery = store.FindPagingQueryRequest

// ProductFindByCaseIdQuery 按案件ID分页查询命令
type ProductFindByCaseIdQuery struct {
	store.FindPagingQueryRequest
	CaseId string `json:"caseId" query:"case-id"`
}

func (q *ProductFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("case_id=='%s'", q.CaseId), nil
}

// ProductFindByTagIdQuery 按案件ID和标签ID查询命令
type ProductFindByTagIdQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	TagId  string `json:"tagId" query:"tag-id" validate:"required" title:"标签ID"`
}

// ProductFindTotalQuery 统计查询命令
type ProductFindTotalQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	Icon   string `json:"icon" query:"icon" title:"图标"`
}

// ProductFindTotalResult 统计结果
type ProductFindTotalResult struct {
	Count int64  `json:"count" title:"产品总数"`
	Icon  string `json:"icon" title:"图标"`
}