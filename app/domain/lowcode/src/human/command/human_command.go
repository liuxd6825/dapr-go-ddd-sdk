package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type HumanCreateCommand struct {
	xbase.BaseCommand
	Data HumanCreateData `json:"data" bson:"data"`
}

type HumanCreateData struct {
	Id   string `json:"id" validate:"required" title:"Id"`
	Code string `json:"code" validate:"required" title:"编号"`
	Name string `json:"name" validate:"required" title:"名称"`
}

type HumanUpdateCommand struct {
	xbase.BaseCommand
	Data HumanUpdateData `json:"data" bson:"data"`
}

type HumanUpdateData struct {
	Id   string `json:"id" validate:"required" title:"Id"`
	Name string `json:"name" validate:"required" title:"名称"`
}
