package ddd_repository

import "fmt"

type FindAutoCompleteQuery interface {
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

	GetField() string
	SetField(string)

	GetValue() string
	SetValue(string)

	GetMustWhere() string
}

type FindByCaseId interface {
	GetCaseId() string
	SetCaseId(string)
}

type FindAutoCompleteQueryRequest struct {
	TenantId string `json:"tenantId" bson:"tenant_id" gorm:"tenant_id"`
	CaseId   string `json:"caseId" bson:"case_id" gorm:"case_id"`
	Field    string `json:"field" bson:"field" gorm:"field"`
	Value    string `json:"value" bson:"value" gorm:"value"`
	PageNum  int64  `json:"pageNum" bson:"page_num" gorm:"page_num"`
	PageSize int64  `json:"pageSize" bson:"page_size" gorm:"page_size"`
	Sort     string `json:"sort" bson:"sort" gorm:"sort"`
	Fields   string `json:"fields" bson:"fields" gorm:"fields"`
	Filter   string `json:"filter" bson:"filter" gorm:"filter"`
}

type FindAutoCompleteQueryDTO = FindAutoCompleteQueryRequest

func NewFindAutoCompleteQuery() FindAutoCompleteQuery {
	return &FindAutoCompleteQueryRequest{}
}

func NewFindAutoCompleteQueryDTO() *FindAutoCompleteQueryDTO {
	return &FindAutoCompleteQueryDTO{}
}

func (f *FindAutoCompleteQueryRequest) GetField() string {
	return f.Field
}

func (f *FindAutoCompleteQueryRequest) SetField(s string) {
	f.Field = s
}

func (f *FindAutoCompleteQueryRequest) GetValue() string {
	return f.Value
}

func (f *FindAutoCompleteQueryRequest) SetValue(s string) {
	f.Value = s
}

func (f *FindAutoCompleteQueryRequest) GetMustWhere() string {
	if len(f.CaseId) > 0 && len(f.Field) > 0 {
		return fmt.Sprintf(`caseId=="%s" and %s~="%s"`, f.CaseId, f.Field, f.Value)
	}
	if len(f.Field) > 0 {
		return fmt.Sprintf(`%s~="%s"`, f.Field, f.Field)
	}
	return ""
}

func (f *FindAutoCompleteQueryRequest) GetTenantId() string {
	return f.TenantId
}

func (f *FindAutoCompleteQueryRequest) SetTenantId(s string) {
	f.TenantId = s
}

func (f *FindAutoCompleteQueryRequest) GetCaseId() string {
	return f.CaseId
}

func (f *FindAutoCompleteQueryRequest) SetCaseId(s string) {
	f.CaseId = s
}

func (f *FindAutoCompleteQueryRequest) GetFields() string {
	return f.Fields
}

func (f *FindAutoCompleteQueryRequest) SetFields(s string) {
	f.Fields = s
}

func (f *FindAutoCompleteQueryRequest) GetFilter() string {
	return f.Filter
}

func (f *FindAutoCompleteQueryRequest) SetFilter(s string) {
	f.Filter = s
}

func (f *FindAutoCompleteQueryRequest) GetSort() string {
	return f.Sort
}

func (f *FindAutoCompleteQueryRequest) SetSort(s string) {
	f.Sort = s
}

func (f *FindAutoCompleteQueryRequest) GetPageNum() int64 {
	return f.PageNum
}

func (f *FindAutoCompleteQueryRequest) SetPageNum(i int64) {
	f.PageNum = i
}

func (f *FindAutoCompleteQueryRequest) GetPageSize() int64 {
	return f.PageSize
}

func (f *FindAutoCompleteQueryRequest) SetPageSize(i int64) {
	f.PageSize = i
}

func (f *FindAutoCompleteQueryRequest) GetQuery() FindAutoCompleteQuery {
	return f
}
