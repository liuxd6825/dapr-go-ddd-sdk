package ppt

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/docx"

func NewReader() *docx.Reader {
	return docx.NewReader()
}
