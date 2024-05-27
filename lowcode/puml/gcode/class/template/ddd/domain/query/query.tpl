package query

import (
	"github.com/liuxu6825/duxm-graph-service/pkg/query-service/domain/graph/view"
)

// GraphFindByIdQuery 按ID查询命令
type GraphFindByIdQuery struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

// GraphFindByIdsQuery 按多个ID查询命令
type GraphFindByIdsQuery struct {
	TenantId string   `json:"tenantId"`
	Ids      []string `json:"ids"`
}

// GraphFindAllQuery 查询所有命令
type GraphFindAllQuery struct {
	TenantId string `json:"tenantId"`
}

// GraphFindPagingQuery 分页查询命令
type GraphFindPagingQuery struct {
	TenantId    string `json:"tenantId"`    // 租户id
	Fields      string `json:"fields"`      // 字段值，多个用逗号分隔
	Filter      string `json:"filter"`      // RSQL过滤条件
	Sort        string `json:"sort"`        // 排序条件
	PageNum     int64  `json:"pageNum"`     // 当前页号
	PageSize    int64  `json:"pageSize"`    // 页大小
	IsTotalRows bool   `json:"isTotalRows"` // 是否汇总记录数
}

// GraphFindPagingResult 分页查询结果
type GraphFindPagingResult struct {
	Data        []*view.GraphView `json:"data"`                 // 分页数据列表
	TotalRows   *int64            `json:"totalRows,omitempty"`  // 总记录数
	TotalPages  *int64            `json:"totalPages,omitempty"` // 总页数
	PageNum     int64             `json:"pageNum"`              // 当前页号
	PageSize    int64             `json:"pageSize"`             // 页大小
	Filter      string            `json:"filter"`               // RSQL过滤条件
	Fields      string            `json:"fields"`               // 字段值，多个用逗号分隔
	Sort        string            `json:"sort"`                 // 排序条件
	Error       error             `json:"error"`                // 错误
	IsFound     bool              `json:"isFound"`              // 是否找到数据
	IsTotalRows bool              `json:"isTotalRows"`          // 是否统计总记录数
}

func NewGraphFindPagingResult() *GraphFindPagingResult {
	return &GraphFindPagingResult{}
}

func NewGraphFindByIdQuery(tenantId, id string) *GraphFindByIdQuery {
	return &GraphFindByIdQuery{
		TenantId: tenantId,
		Id:       id,
	}
}

func NewGraphFindByIdsQuery(tenantId string, ids []string) *GraphFindByIdsQuery {
	return &GraphFindByIdsQuery{
		TenantId: tenantId,
		Ids:      ids,
	}
}

func NewGraphFindAllQuery(tenantId string) *GraphFindAllQuery {
	return &GraphFindAllQuery{
		TenantId: tenantId,
	}
}

func NewGraphFindPagingQuery(tenantId string, fields string, filter string, sort string, pageNum int64, pageSize int64) *GraphFindPagingQuery {
	return &GraphFindPagingQuery{
		TenantId:    tenantId,
		Fields:      fields,
		Filter:      filter,
		Sort:        sort,
		PageNum:     pageNum,
		PageSize:    pageSize,
		IsTotalRows: true,
	}
}

func (q *GraphFindPagingQuery) GetTenantId() string {
	return q.TenantId
}

func (q *GraphFindPagingQuery) SetTenantId(v string) {
	q.TenantId = v
}

func (q *GraphFindPagingQuery) GetFields() string {
	return q.Fields
}

func (q *GraphFindPagingQuery) SetFields(v string) {
	q.Fields = v
}

func (q *GraphFindPagingQuery) GetFilter() string {
	return q.Filter
}

func (q *GraphFindPagingQuery) SetFilter(v string) {
	q.Filter = v
}

func (q *GraphFindPagingQuery) GetSort() string {
	return q.Sort
}

func (q *GraphFindPagingQuery) SetSort(v string) {
	q.Sort = v
}

func (q *GraphFindPagingQuery) GetPageNum() int64 {
	return q.PageNum
}

func (q *GraphFindPagingQuery) SetPageNum(v int64) {
	q.PageNum = v
}

func (q *GraphFindPagingQuery) GetPageSize() int64 {
	return q.PageSize
}

func (q *GraphFindPagingQuery) SetPageSize(v int64) {
	q.PageSize = v
}

func (q *GraphFindPagingQuery) GetIsTotalRows() bool {
	return q.IsTotalRows
}

func (q *GraphFindPagingQuery) SetIsTotalRows(v bool) {
	q.IsTotalRows = v
}
