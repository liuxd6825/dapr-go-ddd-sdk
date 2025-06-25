package doc_extract

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/docx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/pdf"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/ppt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/txt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/xlsx"
	"github.com/spf13/afero"
	"path"
	"strings"
)

type Reader interface {
	ReadFile(fs afero.Fs, fileName string) (string, error)
}

type Extract struct {
	readers map[string]Reader
}

func NewExtract() *Extract {
	e := &Extract{
		readers: make(map[string]Reader),
	}
	txtReader := txt.NewReader()
	docxReader := docx.NewReader()
	xlsxReader := xlsx.NewReader()
	pdfReader := pdf.NewReader()
	pptReader := ppt.NewReader()
	e.AddExt(".txt", txtReader)
	e.AddExt(".json", txtReader)
	e.AddExt(".md", txtReader)
	e.AddExt(".yaml", txtReader)
	e.AddExt(".xml", txtReader)
	e.AddExt(".js", txtReader)
	e.AddExt(".cs", txtReader)

	e.AddExt(".docx", docxReader)
	e.AddExt(".xlsx", xlsxReader)
	e.AddExt(".ppt", pptReader)
	e.AddExt(".pdf", pdfReader)

	return e
}

func (e *Extract) AddExt(extName string, reader Reader) error {
	extName = strings.ToLower(extName)
	e.readers[extName] = reader
	return nil
}

func (e *Extract) Extract(fs afero.Fs, filename string) (string, error) {
	extName := strings.ToLower(path.Ext(filename))
	if reader, ok := e.readers[extName]; ok {
		return reader.ReadFile(fs, filename)
	}
	return "", fmt.Errorf("unsupported file type: %s", extName)
}

func (e *Extract) IsSupport(filename string) bool {
	extName := strings.ToLower(path.Ext(filename))
	_, ok := e.readers[extName]
	return ok
}
