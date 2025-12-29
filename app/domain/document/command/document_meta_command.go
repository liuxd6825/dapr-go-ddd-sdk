package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type DocumentMetaCreateCommand struct {
	xbase.Command[model.DocumentMeta]
}

type DocumentMetaSubmitCommand struct {
	xbase.Command[model.DocumentMeta]
}

type DocumentMetaCreateManyCommand struct {
	xbase.Command[[]*model.DocumentMeta]
}

type DocumentMetaUpdateCommand struct {
	xbase.Command[model.DocumentMeta]
}

type DocumentMetaSaveStatusBySourceType struct {
	xbase.Command[model.DocumentMeta]
}

type DocumentMetaSaveStatusBySourceTypeData struct {
	Id         string
	CaseId     string
	DocumentId string
	SourceType string
	Source     string
	Name       string
	Value      string
}
