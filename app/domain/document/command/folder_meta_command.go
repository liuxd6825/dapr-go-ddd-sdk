package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type FolderMetaCreateCommand struct {
	xbase.Command[model.FolderMeta]
}

type FolderMetaCreateManyCommand struct {
	xbase.Command[[]*model.FolderMeta]
}

type FolderMetaUpdateCommand struct {
	xbase.Command[model.FolderMeta]
}
