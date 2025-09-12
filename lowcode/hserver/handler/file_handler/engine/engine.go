package engine

import (
	"github.com/flosch/pongo2/v6"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/handler/file_handler/common"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
)

type Engine struct {
	serverFs    afero.Fs
	webFs       afero.Fs
	loader      *common.Loader         // Afero 文件加载器
	extension   string                 // 模板文件扩展名
	tplSet      *pongo2.TemplateSet    // Pongo2 模板集合
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
func NewEngine(env *env.Env, serverFs afero.Fs, webFs afero.Fs, extension string) *Engine {
	loader := common.NewLoader(webFs)
	tplSet := pongo2.NewSet("afero", loader)
	engine := &Engine{
		serverFs:    serverFs,
		webFs:       webFs,
		loader:      loader,
		extension:   extension,
		env:         env,
		tplSet:      tplSet,
		templateMap: types.NewCMap[*pongo2.Template](),
		funcs:       make(map[string]interface{}),
	}

	schema := NewSchemaTemplate(tplSet, serverFs, webFs)
	sheet := NewSheetTemplate(tplSet, serverFs, webFs)
	grid := NewAgGridTemplate(tplSet, serverFs, webFs)
	form := NewFormTemplate(tplSet, serverFs, webFs)
	query := NewQueryTemplate(tplSet, serverFs, webFs)
	include := NewInclude(engine, serverFs, webFs)
	// 在注册时
	//set.Globals["include"] = engin.includeHTML()
	//engine.AddFunc("litSSR", litSSR)
	engine.AddFunc("schema", schema.Render)
	engine.AddFunc("sheet", sheet.Render)
	engine.AddFunc("grid", grid.Render)
	engine.AddFunc("form", form.Render)
	engine.AddFunc("query", query.Render)
	engine.AddFunc("include", include.Render)
	engine.AddFunc("ifElse", IfElse)
	engine.AddFunc("toJsonString", ToJsonString)
	engine.AddFunc("nullQuery", NullQuery)
	engine.AddFunc("nullForm", NullForm)
	engine.AddFunc("formValue", FormValue)
	engine.AddFunc("nullColumn", NullColumn)
	engine.AddFunc("columnValue", ColumnValue)
	engine.AddFunc("onlyField", OnlyField)
	engine.AddFunc("onlyOneField", OnlyOneField)
	engine.AddFunc("mapValue", MapValue)
	engine.AddFunc("propValue", PropValue)
	engine.AddFunc("global", GetGlobal)
	return engine
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

// Load 加载模板（无需操作，Afero 动态加载）
func (e *Engine) Load() error {
	return nil
}

// ExecuteWriter 渲染模板
func (e *Engine) ExecuteWriter(w io.Writer, filename string, layout string, bindingData any) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			err = errors.NewErr(err, "file_handler.Engine.ExecuteWriter()")
		}
	}()
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
		fileContent, err = afero.ReadFile(e.webFs, fileName)
		if err != nil {
			return err
		}
	}

	if tmpl == nil {
		// 获取模板
		tmpl, err = e.tplSet.FromBytes(fileContent)
		if err != nil {
			return err
		}
		if prodMode {
			e.templateMap.Set(fileName, tmpl)
		}
	}

	if _, ok = w.(iris.Context); !ok {
		return errors.New("the view engine binding context is not of type iris.Context")
	}

	data["_workDir"] = filepath.Dir(fileName)

	// 渲染模板
	return tmpl.ExecuteWriter(data, w)
}

// AddFunc 添加自定义模板函数
func (e *Engine) AddFunc(name string, fn interface{}) {
	e.funcs[name] = fn
	e.tplSet.Globals[name] = fn
}
