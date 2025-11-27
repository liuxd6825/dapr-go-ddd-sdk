package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type DocumentMetaCreateCommand struct {
	xbase.Command[model.DocumentMeta]
}

type DocumentMetaCreateManyCommand struct {
	xbase.Command[[]*model.DocumentMeta]
}

type DocumentMetaUpdateCommand struct {
	xbase.Command[model.DocumentMeta]
}
