package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanReportedAmountFindByIdQuery 按ID查询
type HumanReportedAmountFindByIdQuery = store.FindByIdRequest

// HumanReportedAmountFindPagingQuery 分页查询
type HumanReportedAmountFindPagingQuery = store.FindPagingQueryRequest
