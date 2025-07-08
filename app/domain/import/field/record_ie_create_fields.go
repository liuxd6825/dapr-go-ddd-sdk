package field

type RecordIeCreateFields struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
	Field    string `json:"field"`
	Value    any    `json:"value"`
}

type RecordIeUpdateFieldFields struct {
	TenantId string `json:"tenantId"`

	Id     string         `json:"id"`
	Values map[string]any `json:"values"`
}

type RecordIeUpdateFilterFields struct {
	TenantId string         `json:"tenantId"`
	Filter   string         `json:"filter"`
	TaskId   string         `json:"taskId"`
	Values   map[string]any `json:"values"`
}
