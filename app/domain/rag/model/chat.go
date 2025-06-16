package model

type Chat struct {
	Id       string `gorm:"id;primaryKey" json:"id"`
	TenantId string `gorm:"tenant_id" json:"tenantId,omitempty"` // 租户ID
	CaseId   string `gorm:"case_id" json:"caseId,omitempty"`
	Title    string `gorm:"title" json:"title,omitempty"` // 标题
}
