package field

type RecordIeDeleteFields struct {
	TenantId string `json:"tenantId"`
	TaskId   string `json:"taskId"`
	Id       string `json:"id,omitempty"  desc:"租户标识"` // 行Id
}
