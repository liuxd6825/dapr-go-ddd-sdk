package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"

// RecordDayFindByIdQuery 按ID查询命令
type RecordDayFindByIdQuery = ddd_query.FindByIdQuery

type RecordDayFindBySumChartQuery struct {
	CaseId      string
	SummaryType string
	Filter      string
}

type RecordDayFindBySumTableQuery struct {
	CaseId      string
	Filter      string
	GroupFilter string
	Sort        string
	PageSize    int64
	PageNum     int64
}

// RecordDayFindByCaseIdQuery 按聚合根ID查询命令
type RecordDayFindByCaseIdQuery struct {
	ddd_query.FindPagingQuery
	CaseId string `json:"caseId" query:"case-id"`
}
