package field

import (
	"time"
)

//
// AccountNodeCreateFields
// @Description: 添加节点命令
//
type AccountNodeCreateFields struct {
	Id              string            `json:"id" validate:"required" `       // 节点Id
	CaseId          string            `json:"caseId" validate:"required" `   // 案件ID
	GraphId         string            `json:"graphId" validate:"required" `  // 聚合根ID
	TenantId        string            `json:"tenantId" validate:"required" ` // 租户ID
	Labels          []string          `json:"labels"`                        // 标签
	Account         string            `json:"account"`                       // 账号
	AccountTypeName string            `json:"accountTypeName"`               // 账号类型
	AccountTypeId   string            `json:"accountTypeId"`                 // 账号类型Id
	AccountUsers    any               `json:"accountUsers"`                  // 使用人信息
	BankId          string            `json:"bankId"`                        // 银行Id
	BankName        string            `json:"bankName"`                      // 银行
	Category        Category          `json:"category"`                      // 类别(1，个人；2，公司)
	LogoutTime      *time.Time        `json:"logoutTime"`                    // 注销时间
	OpenTime        *time.Time        `json:"openTime"`                      // 开户时间
	Owner           string            `json:"owner"`                         // 所有者
	OwnerId         string            `json:"ownerId"`                       // 所有者Id
	OwnerIdent      string            `json:"ownerIdent"`                    // 所有者证件号
	Status          AccountNodeStatus `json:"status"`                        // 状态（1，正常；0，已冻结；-1，已注销；-2，异常）
	Remarks         string            `json:"remarks"`                       // 备注
	IsDeleted       bool              `json:"isDeleted"`                     // 是否删除
}

func (f *AccountNodeCreateFields) GetId() string {
	return f.Id
}

func (f *AccountNodeCreateFields) SetId(v string) {
	f.Id = v
}

func (f *AccountNodeCreateFields) GetTenantId() string {
	return f.TenantId
}

func (f *AccountNodeCreateFields) SetTenantId(v string) {
	f.TenantId = v
}
