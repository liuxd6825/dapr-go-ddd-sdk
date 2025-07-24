package command

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"

type TagRelationCreateCommand struct {
	CommandId string            `json:"commandId"`
	Data      model.TagRelation `json:"data"`
}

type TagRelationUpdateCommand struct {
	CommandId string            `json:"commandId"`
	Data      model.TagRelation `json:"data"`
}

type TagRelationDeleteCommand struct {
	CommandId string            `json:"commandId"`
	Data      model.TagRelation `json:"data"`
}
