package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type TagRelationCreateCommand struct {
	xbase.Command[model.TagRelation]
}

type TagRelationUpdateCommand struct {
	xbase.Command[model.TagRelation]
}

type TagRelationDeleteCommand struct {
	xbase.DeleteByIdCommand
}
