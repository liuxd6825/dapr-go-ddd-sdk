package model

//
// AccountNode
// @Description: 图节点 实体类型
//
type AccountNode struct {
	CaseId   string    `json:"caseId"  validate:"required"` // 备注
	GraphId  string    `json:"graphId"  validate:"gt=0"`
	Id       string    `json:"id"  validate:"required"`       // 节点Id
	Labels   []*string `json:"labels"  validate:"-"`          // 标签
	Name     string    `json:"name"  validate:"-"`            // 节点名称
	Remarks  string    `json:"remarks"  validate:"-"`         // 备注
	TenantId string    `json:"tenantId"  validate:"required"` // 租户ID
}

//
// NewAccountNode
// @Description: 新建图节点对象
//
func NewAccountNode() *AccountNode {
	return &AccountNode{}
}

//
// GetId
// @Description: 取ID值
//
func (e *AccountNode) GetId() string {
	return e.Id
}

func (e *AccountNode) SetId(v string) {
	e.Id = v
}
