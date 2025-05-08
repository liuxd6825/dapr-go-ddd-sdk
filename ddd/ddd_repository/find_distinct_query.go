package ddd_repository

import "fmt"

type FindDistinctQuery interface {
	GetTenantId() string
	SetTenantId(string)

	GetFields() string
	SetFields(string)

	GetFilter() string
	SetFilter(string)

	GetSort() string
	SetSort(string)

	GetPageNum() int64
	SetPageNum(int64)

	GetPageSize() int64
	SetPageSize(int64)

	GetCaseId() string
	SetCaseId(string)

	GetGroupCols() []*GroupCol
	SetGroupCols([]*GroupCol)

	GetMustWhere() string
}

type FindDistinctQueryRequest struct {
	TenantId  string      `json:"tenantId" bson:"tenant_id" gorm:"tenant_id"`
	CaseId    string      `json:"caseId" bson:"case_id" gorm:"case_id"`
	PageNum   int64       `json:"pageNum" bson:"page_num" gorm:"page_num"`
	PageSize  int64       `json:"pageSize" bson:"page_size" gorm:"page_size"`
	Sort      string      `json:"sort" bson:"sort" gorm:"sort"`
	Fields    string      `json:"fields" bson:"fields" gorm:"fields"`
	Filter    string      `json:"filter" bson:"filter" gorm:"filter"`
	GroupCols []*GroupCol `json:"groupCols" bson:"group_cols" gorm:"group_cols"`
}

type FindDistinctQueryDTO struct {
	TenantId  string `json:"tenantId" bson:"tenant_id" gorm:"tenant_id"`
	CaseId    string `json:"caseId" bson:"case_id" gorm:"case_id"`
	PageNum   int64  `json:"pageNum" bson:"page_num" gorm:"page_num"`
	PageSize  int64  `json:"pageSize" bson:"page_size" gorm:"page_size"`
	Sort      string `json:"sort" bson:"sort" gorm:"sort"`
	Fields    string `json:"fields" bson:"fields" gorm:"fields"`
	Filter    string `json:"filter" bson:"filter" gorm:"filter"`
	GroupCols string `json:"groupCols" bson:"group_cols" gorm:"group_cols"`
}

func NewFindDistinctQuery() FindDistinctQuery {
	return &FindDistinctQueryRequest{}
}

func NewFindDistinctQueryDTO() *FindDistinctQueryDTO {
	return &FindDistinctQueryDTO{}
}

func (f *FindDistinctQueryDTO) GetQuery() FindDistinctQuery {
	r := &FindDistinctQueryRequest{}
	r.TenantId = f.TenantId
	r.CaseId = f.CaseId

	r.PageNum = f.PageNum
	r.PageSize = f.PageSize
	r.Sort = f.Sort

	r.Filter = f.Filter
	r.Fields = f.Fields
	r.GroupCols = NewGroupCols(f.GroupCols).Cols()
	return r
}

func (f *FindDistinctQueryRequest) GetMustWhere() string {
	if len(f.CaseId) > 0 {
		return fmt.Sprintf(`caseId=="%s"`, f.CaseId)
	}
	return ""
}

func (f *FindDistinctQueryRequest) GetTenantId() string {
	return f.TenantId
}

func (f *FindDistinctQueryRequest) SetTenantId(s string) {
	f.TenantId = s
}

func (f *FindDistinctQueryRequest) GetCaseId() string {
	return f.CaseId
}

func (f *FindDistinctQueryRequest) SetCaseId(s string) {
	f.CaseId = s
}

func (f *FindDistinctQueryRequest) GetFields() string {
	return f.Fields
}

func (f *FindDistinctQueryRequest) SetFields(s string) {
	f.Fields = s
}

func (f *FindDistinctQueryRequest) GetFilter() string {
	return f.Filter
}

func (f *FindDistinctQueryRequest) SetFilter(s string) {
	f.Filter = s
}

func (f *FindDistinctQueryRequest) GetSort() string {
	return f.Sort
}

func (f *FindDistinctQueryRequest) SetSort(s string) {
	f.Sort = s
}

func (f *FindDistinctQueryRequest) GetPageNum() int64 {
	return f.PageNum
}

func (f *FindDistinctQueryRequest) SetPageNum(i int64) {
	f.PageNum = i
}

func (f *FindDistinctQueryRequest) GetPageSize() int64 {
	return f.PageSize
}

func (f *FindDistinctQueryRequest) SetPageSize(i int64) {
	f.PageSize = i
}

func (f *FindDistinctQueryRequest) GetGroupCols() []*GroupCol {
	return f.GroupCols
}

func (f *FindDistinctQueryRequest) SetGroupCols(cols []*GroupCol) {
	f.GroupCols = cols
}
