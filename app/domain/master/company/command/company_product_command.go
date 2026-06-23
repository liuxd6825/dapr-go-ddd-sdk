package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)


type CompanyProductCreateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.CompanyProduct `json:"data"`
}

type CompanyProductUpdateCommand struct {
	CommandId string                  `json:"commandId"`
	Data      []*model.CompanyProduct `json:"data"`
}

type CompanyProductSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.CompanyProduct `json:"insertData"`
		UpdateData []*model.CompanyProduct `json:"updateData"`
	} `json:"data"`
}