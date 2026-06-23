package query

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// ContractFindByIdQuery 按ID查询命令
type ContractFindByIdQuery = store.FindByIdRequest

// ContractFindPagingQuery 分页查询命令
type ContractFindPagingQuery = store.FindPagingQueryRequest

// ContractFindByCaseIdQuery 按案件ID分页查询命令
type ContractFindByCaseIdQuery struct {
	store.FindPagingQueryRequest
	CaseId string `json:"caseId" query:"case-id"`
}

func (q *ContractFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("case_id=='%s'", q.CaseId), nil
}

// ContractFindByTagIdQuery 按案件ID和标签ID查询命令
type ContractFindByTagIdQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	TagId  string `json:"tagId" query:"tag-id" validate:"required" title:"标签ID"`
}

// ContractFindTotalQuery 统计查询命令
type ContractFindTotalQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	Icon   string `json:"icon" query:"icon" title:"图标"`
}

// ContractFindTotalResult 统计结果
type ContractFindTotalResult struct {
	Count int64  `json:"count" title:"合同总数"`
	Icon  string `json:"icon" title:"图标"`
}