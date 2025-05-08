package transform

import (
	_ "embed" // we need this for embedding Babel
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
)

type Babel struct {
	this      goja.Value
	transform goja.Callable
	vm        *goja.Runtime
}

var defaultOpts = map[string]interface{}{
	"plugins": []any{
		[]any{
			"transform-es2015-modules-commonjs",
			map[string]any{"loose": false},
		},
	},
	"ast":           false,
	"sourceMaps":    false,
	"babelrc":       false,
	"compact":       false,
	"retainLines":   true,
	"highlightCode": false,
}

// 使用 Babel 转换 ES6 -> ES5
const babelTransform = `
		(function(code) {
			return Babel.transform(code, {
				plugins:[["transform-es2015-modules-commonjs", {loose: false}]],
				ast:false,
				sourceMaps:false,
				babelrc:false,
				compact:false,
				retainLines:true,
				highlightCode:false,
			}).code;
		})
	`

//go:embed lib/babel.min.js
var babelSrc string

func NewBabel() (*Babel, error) {
	vm := goja.New()

	_, err := vm.RunString(babelSrc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("加载 Babel 失败： %s", err))
	}

	// 将 Babel 函数注入到 Goja 中
	transformFunc, err := vm.RunString(babelTransform)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("定义 Babel 转换函数失败： %s", err))
	}

	// 转换 ES6 代码为 ES5
	transform, ok := goja.AssertFunction(transformFunc)
	if !ok {
		return nil, errors.New("转换函数断言失败")

	}

	this := vm.Get("Babel")
	return &Babel{
		this:      this,
		transform: transform,
		vm:        vm,
	}, nil
}

func (b *Babel) Transform(es6Code, filename string, sourceMapsEnabled bool, inputSrcMap []byte) ([]byte, error) {
	opts := make(map[string]interface{})
	for k, v := range defaultOpts {
		opts[k] = v
	}
	if sourceMapsEnabled {
		// given that the source map should provide accurate lines(and columns), this option isn't needed
		// it also happens to make very long and awkward lines, especially around import/exports and definitely a lot
		// less readable overall. Hopefully it also has some performance improvement not trying to keep the same lines
		opts["retainLines"] = false
		opts["sourceMaps"] = true
		if inputSrcMap != nil {
			srcMap := new(map[string]interface{})
			if err := json.Unmarshal(inputSrcMap, &srcMap); err != nil {
				return nil, err
			}
			opts["inputSourceMap"] = srcMap
		}
	}
	if filename != "" {
		opts["filename"] = filename
	}

	es5Code, err := b.transform(b.this, b.vm.ToValue(es6Code), b.vm.ToValue(opts))
	if err != nil {
		return nil, err
	}
	return []byte(es5Code.String()), nil
}
