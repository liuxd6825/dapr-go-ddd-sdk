package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)


type CompanyContractCreateCommand struct {
	CommandId string                     `json:"commandId"`
	Data      []*model.CompanyContract   `json:"data"`
}

type CompanyContractUpdateCommand struct {
	CommandId string                     `json:"commandId"`
	Data      []*model.CompanyContract   `json:"data"`
}

type CompanyContractSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.CompanyContract `json:"insertData"`
		UpdateData []*model.CompanyContract `json:"updateData"`
	} `json:"data"`
}