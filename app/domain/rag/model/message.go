package model

type Message struct {
	Id       string `gorm:"id;primaryKey" json:"id"`
	ChatId   string `gorm:"chat_id" json:"chatId"`
	TenantId string `gorm:"tenant_id" json:"tenantId,omitempty"` // 租户ID
	CaseId   string `gorm:"case_id" json:"caseId,omitempty"`
	Content  string `gorm:"title" json:"title,omitempty"` // 内容
	ParentId string `gorm:"parent_id" json:"parentId,omitempty"`
}
