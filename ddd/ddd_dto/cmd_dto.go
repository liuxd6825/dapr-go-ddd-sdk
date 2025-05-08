package ddd_dto

type DeleteByIdRequest struct {
	CommandId string `json:"commandId"`
	Id        string `json:"id"`
}
