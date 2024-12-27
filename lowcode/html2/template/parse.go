package template

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/spf13/afero"
	"strings"
)

type Parse struct {
	repo         string
	ref          string
	htmlPathName string
	htmlFileName string
	content      string
	fs           afero.Fs
	schema       *schema.Schema
	uiSchema     *schema.UiSchema
}

type ParseResult struct {
	HTML         string           `json:"html"`
	Script       strings.Builder  `json:"script"`
	SchemaFile   string           `json:"schemaFile"`
	SchemaName   string           `json:"schemaName"`
	UiSchemaFile string           `json:"uiSchemaFile"`
	UiName       string           `json:"uiName"`
	Template     string           `json:"formTplFile"`
	Schema       *schema.Schema   `json:"schema"`
	UiSchema     *schema.UiSchema `json:"uiSchema"`
}

const (
	SchemaTag   = "schema"
	SlotTag     = "tpl-slot"
	TemplateTag = "template"
	UiSchemaTag = "uischema"
)

func NewParse(fs afero.Fs, htmlFileName string, content []byte) *Parse {
	i := strings.LastIndex(htmlFileName, "/")
	path := htmlFileName[:i]

	return &Parse{
		fs:           fs,
		htmlPathName: path,
		htmlFileName: htmlFileName,
		content:      string(content),
	}
}

func (p *Parse) Parse() (*ParseResult, error) {
	res := &ParseResult{Script: strings.Builder{}}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.content))
	if err != nil {
		return nil, err
	}
	if err := p.initSchema(doc, res); err != nil {
		return nil, err
	}
	if err := p.initBody(doc, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parse) initSchema(doc *goquery.Document, res *ParseResult) error {
	head := doc.Find("head")

	links := head.Find("link")
	links.Each(func(i int, selection *goquery.Selection) {
		id, _ := selection.Attr("id")
		src, _ := selection.Attr("src")
		name, _ := selection.Attr("name")
		switch id {
		case SchemaTag:
			res.SchemaFile = src
			res.SchemaName = name
		case UiSchemaTag:
			res.UiSchemaFile = src
			res.UiName = name
		case TemplateTag:
			res.Template = src
		}
	})

	/*
		if err := p.Valid(res); err != nil {
			return err
		}
	*/

	scripts := head.Find("script")
	scripts.Each(func(i int, selection *goquery.Selection) {
		html := strings.ReplaceAll(selection.Text(), "\n", "")
		html = strings.ReplaceAll(html, "\t", "")
		html = strings.ReplaceAll(html, "  ", "")
		html = strings.TrimSpace(html)
		if html != "" {
			res.Script.WriteString(selection.Text())
		}
	})

	return nil
}

func (p *Parse) Valid(res *ParseResult) error {
	errs := errors.NewErrors()
	if len(res.UiName) == 0 {
		errs.AddString("uiSchema name is empty")
	}
	if len(res.UiSchemaFile) == 0 {
		errs.AddString("uiSchema src is empty")
	}
	if len(res.SchemaFile) == 0 {
		errs.AddString("schema src is empty")
	}
	if len(res.SchemaName) == 0 {
		errs.AddString("schema name is empty")
	}
	if len(res.Template) == 0 {
		errs.AddString("template is empty")
	}
	if errs.HasError() {
		return errs
	}
	return nil
}

func (p *Parse) initBody(doc *goquery.Document, res *ParseResult) error {
	errs := errors.NewErrors()

	if res.SchemaFile != "" {
		sm, err := loadSchema(p.fs, p.htmlPathName, res.SchemaFile, res.SchemaName)
		if err != nil {
			errs.AddError(newError(res.SchemaFile, err))
		}
		p.schema = sm
	}

	if res.UiSchemaFile != "" {
		ui, err := loadUiSchema(p.fs, p.schema, p.htmlPathName, res.UiSchemaFile, res.UiName)
		if err != nil {
			errs.AddError(newError(res.UiSchemaFile, err))
		}
		p.uiSchema = ui
	}

	if errs.HasError() {
		return errs.NewError()
	}

	body := doc.Find("body")
	if body.Children().Length() == 0 {
		html, err := p.schemaRender(res.Template, p.schema, p.uiSchema)
		if err != nil {
			errs.AddError(newError(res.Template, err))
		}
		res.HTML = html
	} else {
		slot := body.Find(SlotTag)
		if slot.Length() > 0 {
			html, err := p.schemaRender(res.Template, p.schema, p.uiSchema)
			if err != nil {
				errs.AddError(newError(res.Template, err))
			}
			slot.SetHtml(html)
		}
		slot.BeforeHtml("<div></div>")
		html, err := body.Html()
		if err != nil {
			errs.AddError(newError("body.Html()", err))
		}
		res.HTML = html
	}
	return errs.NewError()
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
func (p *Parse) schemaRender(htmlTplFileName string, schema *schema.Schema, uiSchema *schema.UiSchema) (string, error) {
	if schema == nil || uiSchema == nil {
		return "", errors.ErrorOf("param schema or uiSchema is nil")
	}
	builder := NewUiBuilder(schema, uiSchema)
	content, err := readFile(p.fs, htmlTplFileName, p.htmlPathName)
	if err != nil {
		return "", err
	}
	writer := bytes.Buffer{}
	err = builder.BuildBytes(content, &writer)
	return writer.String(), err
}

func (r *ParseResult) Valid() error {
	errs := errors.NewErrors()
	if len(r.UiName) == 0 {
		errs.AddString("UiName is empty.")
	}
	if len(r.UiSchemaFile) == 0 {
		errs.AddString("UiSchemaFile is empty.")
	}
	if len(r.SchemaFile) == 0 {
		errs.AddString("SchemaFile is empty.")
	}
	if len(r.SchemaName) == 0 {
		errs.AddString("SchemaName is empty.")
	}
	if len(r.Template) == 0 {
		errs.AddString("Template is empty.")
	}
	if errs.HasError() {
		return errs
	}
	return nil
}

func newError(s string, err error) error {
	return errors.New(fmt.Sprintf("%s %s", s, err.Error()))
}
