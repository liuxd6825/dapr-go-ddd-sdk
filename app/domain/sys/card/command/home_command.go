package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type HomeCreateCommand struct {
	xbase.Command[model.Home]
}
type HomeUpdateCommand struct {
	xbase.Command[model.Home]
}
type HomeDeleteCommand struct {
	xbase.DeleteByIdCommand
}

type HomeSubmitCommand struct {
	xbase.Command[HomeSubmitCommandData]
}

type HomeSubmitCommandData struct {
	Home  *model.Home                     `json:"home"`
	Group SubmitCommandData[*model.Group] `json:"group"`
	Card  SubmitCommandData[*model.Card]  `json:"card"`
}

type SubmitCommandData[T any] struct {
	I []T      `json:"i"`
	U []T      `json:"u"`
	D []string `json:"d"`
}
