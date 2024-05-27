package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

//
// GraphAggregate
// @Description:  <no value> 聚合类型
//
type {{.Class.Name}}Aggregate struct {
    {{- range $field := .Class.Fields }}
    {{$field.Name}} {{$field.DataType}} `json:"{{$field.FN}}"`   // {{$method.CommentsText}}
    {{- end}}
	IsDeleted bool   `json:"isDeleted" validate:"-"` // 已删除
	Remarks   string `json:"remarks" validate:"-"`   // 备注
}

const AggregateType = "{{.Class.Name}}Aggregate"

// MaskMapper时不复制的属性
var aggMapperRemove []string

func init() {
}

//
// NewGraphAggregate
// @Description: 新建<no value> 聚合对象
// @return *GraphAggregate
//
func New{{.Class.Name}}Aggregate() *{{.Class.Name}}Aggregate {
	return &{{.Class.Name}}Aggregate{
	}
}

//
// NewAggregate
// @Description: 新建当前包中的聚合对象，当前包中只能有一个聚合类型
// @return ddd.Aggregate
//
func NewAggregate() ddd.Aggregate {
	return New{{.Class.Name}}Aggregate()
}

//
// GetAggregateVersion
// @Description: 聚合的版本号
// @receiver a
// @return string 版本号
//
func (a *{{.Class.Name}}Aggregate) GetAggregateVersion() string {
	return "v1.0"
}

//
// GetAggregateType
// @Description: 获取聚合的类型
// @receiver a
// @return string 聚合的类型
//
func (a *{{.Class.Name}}Aggregate) GetAggregateType() string {
	return AggregateType
}

//
// GetAggregateId
// @Description: 获取 聚合id
// @receiver a
// @return string 聚合id
//
func (a *{{.Class.Name}}Aggregate) GetAggregateId() string {
	return a.Id
}

//
// GetTenantId
// @Description: 租户id
// @receiver a
// @return string 租户id
//
func (a *{{.Class.Name}}Aggregate) GetTenantId() string {
	return a.TenantId
}
