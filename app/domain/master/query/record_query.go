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

type DistinctAccountByNameQuery struct {
	CaseId     string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	Name       string `json:"name" query:"name"  title:"开户人"`
	MasterType string `json:"masterType" query:"master-type"  title:"主数据类型"`
	MasterId   string `json:"masterName" query:"master-id" title:"主数据Id"`
}

type DistinctAccountResult struct {
	Account   string `json:"account" title:"账号"`
	BankName  string `json:"bankName" title:"银行"`
	OwnerName string `json:"ownerName" title:"开户人"`
}

type DistinctHumanResult struct {
	Name string `json:"name" title:"公司名"`
}

type DistinctCompanyResult struct {
	Name string `json:"name" title:"公司名"`
}

type DistinctNameQuery struct {
	CaseId     string `json:"caseId" query:"case-id" validate:"required" title:"案件ID"`
	MasterType string `json:"masterType" query:"master-type"  title:"主数据类型"`
	MasterId   string `json:"masterName" query:"master-id" title:"主数据Id"`
	MyFilter   string `json:"myFilter" query:"my-filter" title:"账号"`
	OppFilter  string `json:"oppFilter" query:"opp-filter"  title:"银行"`
}

type DistinctNameResult struct {
	Name string `json:"name" title:"公司名"`
}

func (q *RecordFindByCaseIdQuery) GetMustWhere() (string, error) {
	if len(q.CaseId) == 0 {
		return "", errors.New("caseId参数不能为空")
	}
	return fmt.Sprintf("CaseId=='%s'", q.CaseId), nil
}
