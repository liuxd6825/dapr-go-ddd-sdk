package doc_extract

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/docx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/pdf"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/ppt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/txt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/xls"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/doc_extract/xlsx"
	"github.com/sirupsen/logrus"
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

func NewExtract(logger *logrus.Logger) *Extract {
	e := &Extract{
		readers: make(map[string]Reader),
	}
	txtReader := txt.NewReader()
	docxReader := docx.NewReader()
	xlsxReader := xlsx.NewReader()
	xlsReader := xls.NewReader(logger)
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
	e.AddExt(".xls", xlsReader)

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
		text, err := reader.ReadFile(fs, filename)
		if err == nil {
			return e.clear(text), nil
		}
	}
	return "", fmt.Errorf("unsupported file type: %s", extName)
}

func (e *Extract) IsSupport(filename string) bool {
	extName := strings.ToLower(path.Ext(filename))
	_, ok := e.readers[extName]
	return ok
}

func (e *Extract) clear(text string) string {
	text = strings.ReplaceAll(text, "\u0001", " ")
	text = strings.ReplaceAll(text, "\u0014", " ")
	text = strings.ReplaceAll(text, "\u0015", " ")
	text = strings.ReplaceAll(text, "�", " ")
	return text
}
