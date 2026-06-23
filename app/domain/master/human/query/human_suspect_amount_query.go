package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanSuspectAmountFindByIdQuery 按ID查询
type HumanSuspectAmountFindByIdQuery = store.FindByIdRequest

// HumanSuspectAmountFindPagingQuery 分页查询
type HumanSuspectAmountFindPagingQuery = store.FindPagingQueryRequest
