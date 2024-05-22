package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
	"github.com/spf13/afero"
	"strings"
)

type Parse struct {
	repo     string
	ref      string
	fileName string
	content  string
	fs       afero.Fs
	schema   *schema.Schema
	uiSchema *schema.UiSchema
}

type ParseResult struct {
	HTML         string           `json:"html"`
	Script       strings.Builder  `json:"script"`
	SchemaFile   string           `json:"schemaFile"`
	UiSchemaFile string           `json:"uiSchemaFile"`
	FormTplFile  string           `json:"formTplFile"`
	Schema       *schema.Schema   `json:"schema"`
	UiSchema     *schema.UiSchema `json:"uiSchema"`
}

func NewParse(fs afero.Fs, fileName string, content []byte) *Parse {
	return &Parse{
		fs:       fs,
		fileName: fileName,
		content:  string(content),
	}
}

func (h *Parse) Parse() (*ParseResult, error) {
	res := &ParseResult{Script: strings.Builder{}}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(h.content))
	if err != nil {
		return nil, err
	}
	if err := h.initSchema(doc, res); err != nil {
		return nil, err
	}
	if err := h.initBody(doc, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h *Parse) initSchema(doc *goquery.Document, res *ParseResult) error {
	head := doc.Find("head")

	links := head.Find("link")
	links.Each(func(i int, selection *goquery.Selection) {
		id, _ := selection.Attr("id")
		src, _ := selection.Attr("src")
		switch id {
		case "schema":
			res.SchemaFile = src
		case "uischema":
			res.UiSchemaFile = src
		case "formtpl":
			res.FormTplFile = src
		}
	})

	scripts := head.Find("script")
	scripts.Each(func(i int, selection *goquery.Selection) {
		res.Script.WriteString(selection.Text())
	})

	return nil
}

func (h *Parse) initBody(doc *goquery.Document, res *ParseResult) error {
	var sm *schema.Schema
	var ui *schema.UiSchema
	var err error

	if res.SchemaFile != "" {
		sm, err = loadSchema(h.fs, res.SchemaFile)
		if err != nil {
			return newError(res.SchemaFile, err)
		}
		h.schema = sm
	}

	if res.UiSchemaFile != "" {
		ui, err = loadUiSchema(h.fs, sm, res.UiSchemaFile)
		if err != nil {
			return newError(res.UiSchemaFile, err)
		}
		h.uiSchema = ui
	}

	body := doc.Find("body")
	if body.Children().Length() == 0 {
		html, err := h.schemaRender(res.FormTplFile, sm, ui)
		if err != nil {
			return newError(res.FormTplFile, err)
		}
		res.HTML = html
	} else {
		slot := body.Find("hyk-slot")
		if slot.Length() > 0 {
			html, err := h.schemaRender(res.FormTplFile, sm, ui)
			if err != nil {
				return newError(res.FormTplFile, err)
			}
			slot.SetHtml(html)
		}
		html, err := body.Html()
		if err != nil {
			return newError("body.Html()", err)
		}
		res.HTML = html
	}
	return nil
}

// schemaRender
//
//	@Description:  通过uiSchema生成html模板
//	@receiver h
//	@param htmlTplFileName form模板文件名称
//	@param schemaByte  schema文件内容
//	@param uiSchemaByte ui schema文件内容
//	@return string 使用form模板渲染的HTML
//	@return error
func (h *Parse) schemaRender(htmlTplFileName string, sm *schema.Schema, ui *schema.UiSchema) (string, error) {
	builder := NewUiBuilder(sm, ui)
	content, err := readFile(h.fs, htmlTplFileName)
	if err != nil {
		return "", err
	}
	writer := bytes.Buffer{}
	err = builder.BuildBytes(content, &writer)
	return writer.String(), err
}

func newError(s string, err error) error {
	return errors.New(fmt.Sprintf("%s %s", s, err.Error()))
}
