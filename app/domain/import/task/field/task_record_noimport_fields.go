package field

type TaskRecordNoImportFields struct {
	Id       string `json:"id" desc:"Id"`
	TenantId string `json:"tenantId" desc:"租户Id"`
	SchemaId string `json:"schemaId" desc:"主数据类型"`
}
