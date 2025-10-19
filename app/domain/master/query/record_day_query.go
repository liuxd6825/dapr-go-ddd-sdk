package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"

// RecordDayFindByIdQuery 按ID查询命令
type RecordDayFindByIdQuery = ddd_query.FindByIdQuery

// RecordDayFindByCaseIdQuery 按聚合根ID查询命令
type RecordDayFindByCaseIdQuery = ddd_query.FindPagingByCaseIdQuery

type RecordDayFindBySumChartQuery struct {
	CaseId      string `json:"caseId"  param:"caseId" validate:"required" title:"项目ID"`
	SummaryType string `json:"summaryType"  param:"summaryType" validate:"required" title:"汇总方式"`
	Filter      string `json:"filter"  param:"filter" validate:"" title:"过滤条件"`
}

type RecordDayFindBySumTableQuery struct {
	CaseId      string `json:"caseId"  param:"caseId" validate:"required" title:"项目ID"`
	Filter      string `json:"filter"  param:"filter" validate:"" title:"过滤条件"`
	GroupFilter string `json:"groupFilter"  param:"groupFilter"  validate:"" title:"分组条件"`
	Sort        string `json:"sort"  param:"sort"  validate:"" title:"排序"`
	PageSize    int64  `json:"pageSize"  param:"pageSize"  validate:"required" title:"分页大小"`
	PageNum     int64  `json:"pageNum"  param:"pageNum"  validate:""  title:"分页页号"`
}
