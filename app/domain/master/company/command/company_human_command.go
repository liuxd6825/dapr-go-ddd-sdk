package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
)


type CompanyHumanCreateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.CompanyHuman  `json:"data"`
}

type CompanyHumanUpdateCommand struct {
	CommandId string                 `json:"commandId"`
	Data      []*model.CompanyHuman  `json:"data"`
}

type CompanyHumanSubmitManyCommand struct {
	CommandId string `json:"commandId"`
	Data      struct {
		InsertData []*model.CompanyHuman `json:"insertData"`
		UpdateData []*model.CompanyHuman `json:"updateData"`
	} `json:"data"`
}