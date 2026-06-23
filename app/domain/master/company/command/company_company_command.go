package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)


type CompanyCompanyCreateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.CompanyCompany  `json:"data"`
}

type CompanyCompanyUpdateCommand struct {
	CommandId string                   `json:"commandId"`
	Data      []*model.CompanyCompany  `json:"data"`
}

type CompanyCompanySubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.CompanyCompany `json:"insertData"`
		UpdateData []*model.CompanyCompany `json:"updateData"`
	} `json:"data"`
}