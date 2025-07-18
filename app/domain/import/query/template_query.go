package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// TemplateFieldFindPagingQuery 分页查询命令
type TemplateFieldFindPagingQuery struct {
}

// TemplateFieldFindPagingResult 分页查询结果
type TemplateFieldFindPagingResult struct {
	Data        []*model.TemplateField `json:"data"`                 // 分页数据列表
	TotalRows   *int64                 `json:"totalRows,omitempty"`  // 总记录数
	TotalPages  *int64                 `json:"totalPages,omitempty"` // 总页数
	PageNum     int64                  `json:"pageNum"`              // 当前页号
	PageSize    int64                  `json:"pageSize"`             // 页大小
	Filter      string                 `json:"filter"`               // RSQL过滤条件
	Fields      string                 `json:"fields"`               // 字段值，多个用逗号分隔
	Sort        string                 `json:"sort"`                 // 排序条件
	Error       error                  `json:"error"`                // 错误
	IsFound     bool                   `json:"isFound"`              // 是否找到数据
	IsTotalRows bool                   `json:"isTotalRows"`          // 是否统计总记录数
}

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

func NewTemplateFieldFindPagingQuery() *TemplateFieldFindPagingQuery {
	return &TemplateFieldFindPagingQuery{}
}

func NewTemplateFieldFindPagingResult() *TemplateFieldFindPagingResult {
	return &TemplateFieldFindPagingResult{}
}
