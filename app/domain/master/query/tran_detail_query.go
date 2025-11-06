package query

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// TranDetailFindThresholdQuery 查询阈值交易
type TranDetailFindThresholdQuery struct {
	store.FindPagingQueryRequest
	Name      string    `json:"name"  param:"name" validate:"required" title:"账户名称"`
	Amount    float64   `json:"amount"  param:"amount" validate:"required" title:"阈值金额"`
	StartDate time.Time `json:"startDate"  param:"start-date" validate:"required"  title:"开始时间"`
	EndDate   time.Time `json:"endDate"  param:"end-date" validate:"required"  title:"结束时间"`
}

// TranDetailFindCashThresholdQuery 查询现金阈值交易
type TranDetailFindCashThresholdQuery struct {
	store.FindPagingQueryRequest
	Value     float64   `json:"value"  param:"value" validate:"required" title:"阈值"`
	StartDate time.Time `json:"startDate"  param:"start-date" validate:"required"  title:"开始时间"`
	EndDate   time.Time `json:"endDate"  param:"end-date" validate:"required"  title:"结束时间"`
}
