package field

type TaskRecordImportFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
	PageSize int64  `json:"pageSize" desc:"分页大小"`
}
