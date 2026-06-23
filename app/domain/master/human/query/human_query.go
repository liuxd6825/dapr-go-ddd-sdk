package query

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// HumanFindByIdQuery 按ID查询命令
type HumanFindByIdQuery = store.FindByIdRequest

// HumanFindPagingQuery 分页查询命令
type HumanFindPagingQuery = store.FindPagingQueryRequest

// HumanFindByCaseIdQuery 按案件ID分页查询命令
type HumanFindByCaseIdQuery struct {
	store.FindPagingQueryRequest
	CaseId string `json:"caseId" param:"case-id"`
}

func (q *HumanFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("case_id=='%s'", q.CaseId), nil
}

// HumanFindByTagIdQuery 按案件ID和标签ID查询命令
type HumanFindByTagIdQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	TagId  string `json:"tagId" query:"tag-id" validate:"required" title:"标签ID"`
}

// HumanFindTotalQuery 统计查询命令
type HumanFindTotalQuery struct {
	CaseId string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	Icon   string `json:"icon" query:"icon" title:"图标"`
}

// HumanFindTotalResult 统计结果
type HumanFindTotalResult struct {
	Count int64  `json:"count" title:"人员总数"`
	Icon  string `json:"icon" title:"图标"`
}

// PersonTypeItem 人员类型项
type PersonTypeItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// PeopleTypeItem 分析状态项
type PeopleTypeItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
