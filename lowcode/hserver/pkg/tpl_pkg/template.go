package tpl_pkg

import (
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/spf13/afero"
)

var tplCache = cmap.New()

type Template struct {
	cfg    *env.Env
	server element.Server
	fs     afero.Fs
}

func New(cfg *env.Env, server element.Server, fs afero.Fs) *Template {
	return &Template{
		cfg:    cfg,
		server: server,
		fs:     fs,
	}
}

func (e *Template) RenderFile(filename string, data map[string]any) string {
	var txt string
	var err error
	e.getTpl(filename, func(tpl *pongo2.Template) error {
		txt, err = tpl.Execute(data)
		return err
	})
	return txt
}

func (e *Template) RenderString(str string, data map[string]any) string {
	tpl, err := pongo2.FromString(str)
	if err != nil {
		panic(err)
	}
	txt := e.render(tpl, err, data)
	return *txt
}

func (e *Template) renderBytes(bytes []byte, data map[string]any) *string {
	tpl, err := pongo2.FromBytes(bytes)
	if err != nil {
		panic(err)
	}
	return e.render(tpl, err, data)
}

func (e *Template) render(tpl *pongo2.Template, err error, data map[string]any) *string {
	if err != nil {
		panic(err)
	}
	context, err := tpl.Execute(data)
	if err != nil {
		panic(err)
	}
	return &context
}

func (e *Template) getTpl(filename string, fun func(*pongo2.Template) error) {
	tpl := GetCache(filename)
	if tpl != nil {
		if err := fun(tpl); err != nil {
			panic(err)
		}
	}

	bytes, err := fs.ReadFile(e.fs, filename)
	if err != nil {
		panic(err)
	}
	txt := replacePercent(string(bytes))
	tpl, err = pongo2.FromBytes([]byte(txt))
	if err != nil {
		panic(err)
	}
	if err := fun(tpl); err != nil {
		panic(err)
	}
	SetCache(filename, tpl)
}

func GetCache(filename string) *pongo2.Template {
	val, hash := tplCache.Get(filename)
	if hash {
		var tpl = val.(*pongo2.Template)
		return tpl
	}
	return nil
}

func DeleteCache(filename string) {
	tplCache.Remove(filename)
}

func SetCache(filename string, tpl *pongo2.Template) {
	tplCache.Set(filename, tpl)
}

// replacePercent replaces '%' with '%%', except when it appears in '{%' or '%}'
func replacePercent(input string) string {
	runes := []rune(input)
	length := len(runes)
	var result []rune

	for i := 0; i < length; i++ {
		if runes[i] == '%' {
			// Check if this % is part of {% or %}
			isPartOfTag := (i > 0 && runes[i-1] == '{') || (i < length-1 && runes[i+1] == '}')

			if isPartOfTag {
				result = append(result, runes[i])
			} else {
				result = append(result, '%', '%')
			}
		} else {
			result = append(result, runes[i])
		}
	}

	return string(result)
}
