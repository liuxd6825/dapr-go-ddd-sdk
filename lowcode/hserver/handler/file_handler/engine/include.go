package engine

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/spf13/afero"
	"html/template"
)

type Include struct {
	engine   *Engine
	webFs    afero.Fs
	serverFs afero.Fs
}

func NewInclude(engine *Engine, serverFs afero.Fs, webFs afero.Fs) *Include {
	return &Include{engine: engine, serverFs: serverFs, webFs: webFs}
}

func (e *Include) Render(filename string) template.HTML {
	fileName := fileutils.AbsPath(filename, &fileutils.ReadOptions{RootPath: "", WorkPath: ""})
	// 获取文件的绝对路径（可选，根据你的需求调整）
	data, err := afero.ReadFile(e.webFs, fileName)
	if err != nil {
		return template.HTML(fmt.Sprintf("<div> 错误: 无法读取文件 %s: %v </div>", fileName, err))
	}
	tpl, err := e.engine.tplSet.FromBytes(data)
	if err != nil {
		return template.HTML(fmt.Sprintf("<div> 错误: 读取文件渲染时出策 %s: %v </div>", fileName, err))
	}
	html, err := tpl.Execute(nil)
	if err != nil {
		return template.HTML(fmt.Sprintf("<div> 错误: 读取文件渲染时出策 %s: %v </div>", fileName, err))
	}
	return template.HTML(html)
}
