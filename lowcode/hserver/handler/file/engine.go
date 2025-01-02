package file

import (
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
)

type Engine struct {
	fs          afero.Fs
	loader      *Loader                // Afero 文件加载器
	extension   string                 // 模板文件扩展名
	templates   *pongo2.TemplateSet    // Pongo2 模板集合
	funcs       map[string]interface{} // 自定义模板函数
	templateMap *types.CMap[*pongo2.Template]
}

func (e *Engine) Name() string {
	return "aferoPongo2"
}

func (e *Engine) Ext() string {
	return e.extension
}

// NewEngine 创建一个新的 Afero Pongo2 引擎
func NewEngine(fs afero.Fs, extension string) *Engine {
	loader := NewLoader(fs)
	set := pongo2.NewSet("afero", loader)
	return &Engine{
		fs:          fs,
		loader:      loader,
		extension:   extension,
		templates:   set,
		templateMap: types.NewCMap[*pongo2.Template](),
		funcs:       make(map[string]interface{}),
	}
}

// Load 加载模板（无需操作，Afero 动态加载）
func (e *Engine) Load() error {
	return nil
}

// ExecuteWriter 渲染模板
func (e *Engine) ExecuteWriter(w io.Writer, filename string, layout string, bindingData interface{}) (err error) {
	if filepath.Ext(filename) == "" {
		filename += e.extension
	}
	fileName := "/" + filename

	var tmpl *pongo2.Template
	var data []byte
	if t, ok := e.templateMap.Get(fileName); ok {
		tmpl = t
	} else {
		data, err = afero.ReadFile(e.fs, fileName)
		if err != nil {
			return err
		}
	}

	if tmpl == nil {
		// 获取模板
		tmpl, err = e.templates.FromBytes(data)
		if err != nil {
			return err
		}
		e.templateMap.Set(fileName, tmpl)
	}
	// 检查绑定数据
	ctx, ok := bindingData.(map[string]any)
	if !ok {
		return errors.New("binding data should be of type map[string]interface{}")
	}

	// 渲染模板
	return tmpl.ExecuteWriter(ctx, w)
}

// AddFunc 添加自定义模板函数
func (e *Engine) AddFunc(name string, fn interface{}) {
	e.funcs[name] = fn
	e.templates.Globals[name] = fn
}
