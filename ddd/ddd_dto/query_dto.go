package ddd_dto

type FindByIdRequest struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

type FindPagingRequest struct {
	TenantId string `json:"tenantId"`
}
