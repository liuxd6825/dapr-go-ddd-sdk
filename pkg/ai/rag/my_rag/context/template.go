package context

import (
	"strings"
	"text/template"
)

var defaultPromptTemplate = `{{.SystemPrompt}}{{if .Summary}}
背景：{{.Summary}}{{end}}

{{if .History}}
【历史对话】：
{{- range $index, $item := .History}}
{{$item.Role}}: {{$item.Content}}
{{- end}}
{{end}}

{{if .RetrievedContext}}
【参考知识】：
{{- range $index, $chunk := .RetrievedContext}}
[{{add $index 1}}] {{$chunk}}
{{- end}}
{{end}}

问题: {{.CurrentQuery}}

回答：`

type TemplateData struct {
	SystemPrompt     string
	Summary          string
	History          []*Content
	RetrievedContext []string
	CurrentQuery     string
}

func RenderPrompt(data *TemplateData) (string, error) {
	tmpl := template.New("rag-prompt").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	var err error
	tmpl, err = tmpl.Parse(defaultPromptTemplate)
	if err != nil {
		return "", err
	}

	buf := &strings.Builder{}
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderPromptWithTemplate(templ string, data *TemplateData) (string, error) {
	tmpl := template.New("custom-prompt").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	var err error
	tmpl, err = tmpl.Parse(templ)
	if err != nil {
		return "", err
	}

	buf := &strings.Builder{}
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func BuildTemplateData(param *BuildParam, summary string) *TemplateData {
	history := make([]*Content, 0, len(param.ConversationHistory))
	for _, h := range param.ConversationHistory {
		history = append(history, h)
	}

	return &TemplateData{
		SystemPrompt:     "",
		Summary:          summary,
		History:          history,
		RetrievedContext: param.RetrievedContext,
		CurrentQuery:     param.Query,
	}
}