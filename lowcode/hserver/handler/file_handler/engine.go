package file_handler

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/spf13/afero"
	"html/template"
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
	env         *env.Env
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
	engin := &Engine{
		fs:          fs,
		loader:      loader,
		extension:   extension,
		templates:   set,
		templateMap: types.NewCMap[*pongo2.Template](),
		funcs:       make(map[string]interface{}),
	}
	// 在注册时
	//set.Globals["include"] = engin.includeHTML()
	return engin
}

var ctxKey = "iris_ctx"

func RegisterTemplateFunc(app *iris.Application) {
	app.OnAnyErrorCode(func(ctx iris.Context) {
		// 确保错误页面也能获取到Context
		ctx.ViewData(ctxKey, ctx)
		ctx.Next()
	})

	app.Use(func(ctx iris.Context) {
		ctx.ViewData(ctxKey, ctx)
		ctx.Next()
	})
}

func (e *Engine) includeHTML(workDir string) func(string) template.HTML {
	return func(filename string) template.HTML {
		fileName := fileutils.AbsPath(filename, &fileutils.ReadOptions{RootPath: "", WorkPath: workDir})
		// 获取文件的绝对路径（可选，根据你的需求调整）
		data, err := afero.ReadFile(e.fs, fileName)
		if err != nil {
			return template.HTML(fmt.Sprintf("<div> 错误: 无法读取文件 %s: %v </div>", fileName, err))
		}
		return template.HTML(data)
	}
}

// Load 加载模板（无需操作，Afero 动态加载）
func (e *Engine) Load() error {
	return nil
}

// ExecuteWriter 渲染模板
func (e *Engine) ExecuteWriter(w io.Writer, filename string, layout string, bindingData any) (err error) {
	if filepath.Ext(filename) == "" {
		filename += e.extension
	}
	fileName := "/" + filename

	prodMode := e.env.GetProdMode()
	data, ok := bindingData.(map[string]any)
	if !ok {
		return errors.New("the view engine binding data is not of type map[string]any")
	}

	var tmpl *pongo2.Template
	var fileContent []byte
	if !prodMode {
		if fc, ok := data["templateData"]; ok {
			fileContent = fc.([]byte)
		}
	}

	if t, ok := e.templateMap.Get(fileName); ok {
		tmpl = t
	} else if fileContent == nil {
		fileContent, err = afero.ReadFile(e.fs, fileName)
		if err != nil {
			return err
		}
	}

	if tmpl == nil {
		// 获取模板
		tmpl, err = e.templates.FromBytes(fileContent)
		if err != nil {
			return err
		}
		if prodMode {
			e.templateMap.Set(fileName, tmpl)
		}
	}

	ctx, ok := w.(iris.Context)
	if !ok {
		println("ctx:", ctx)
	}

	workDir := filepath.Dir(fileName)
	data["include"] = e.includeHTML(workDir)

	// 渲染模板
	return tmpl.ExecuteWriter(data, w)
}

// AddFunc 添加自定义模板函数
func (e *Engine) AddFunc(name string, fn interface{}) {
	e.funcs[name] = fn
	e.templates.Globals[name] = fn
}
