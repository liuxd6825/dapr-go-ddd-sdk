package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type TemplateFindByIdQuery struct {
	TenantId string `json:"tenantId"` // 租户id
	Id       string `json:"id"`
}

type TemplateFindAllQuery struct {
	TenantId string `json:"tenantId"` // 租户id
}

type TemplateFindQuery struct {
}

// TemplateFindPagingQuery 分页查询命令
type TemplateFindPagingQuery = store.FindPagingQueryRequest

// TemplateFindPagingByCaseIdQuery 分页查询命令
type TemplateFindPagingByCaseIdQuery = ddd_query.FindPagingByCaseIdQuery

// Validate
// @Description: 命令数据验证
func (c *TemplateFindByIdQuery) Validate() error {
	ve := errors.NewVerifyError()
	if len(c.TenantId) == 0 {
		ve.AppendField("TenantId", "不能为空")
	}
	if len(c.Id) == 0 {
		ve.AppendField("BankName", "不能为空")
	}
	return ve.GetError()
}
