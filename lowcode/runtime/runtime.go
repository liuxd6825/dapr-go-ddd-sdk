package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
)

type Runtime struct {
	vm          *goja.Runtime
	moduleCache map[string]*goja.Object
	reader      fs.Reader
}

// 模块缓存

var DefaultPoolSize = 5
var fieldNameMapper = &FieldNameMapper{}

func NewRuntime(reader fs.Reader) *Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(fieldNameMapper)
	r := &Runtime{
		vm:          vm,
		moduleCache: map[string]*goja.Object{},
		reader:      reader,
	}
	r.setRequire(r.vm, reader)
	return r
}

// 注册自定义 require 函数
func (r *Runtime) setRequire(vm *goja.Runtime, fsReader fs.Reader) {
	_ = vm.Set("require", func(call goja.FunctionCall) goja.Value {
		modulePath := call.Argument(0).String()
		value, err := r.require(fsReader, modulePath)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return value
	})
}

// 自定义 require 函数
func (r *Runtime) require(reader fs.Reader, modulePath string) (val goja.Value, err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	moduleName := modulePath
	// 检查缓存
	if cachedModule, ok := r.moduleCache[modulePath]; ok {
		return cachedModule, nil
	}

	// 读取模块文件内容
	content, err := reader.ReadFile(modulePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load module %s: %v", modulePath, err)
	}

	// 创建新的 Runtime
	moduleVM := goja.New()

	// 创建 module 和 exports 对象
	exports := moduleVM.NewObject()
	moduleObject := moduleVM.NewObject()
	_ = moduleObject.Set("exports", exports)

	// 注入 module 和 exports
	_ = moduleVM.Set("module", moduleObject)
	_ = moduleVM.Set("exports", exports)

	// 绑定 require 函数，让模块内可以嵌套调用
	r.setRequire(moduleVM, reader)

	// 包装模块代码，注入 require、module 和 exports
	wrappedCode := fmt.Sprintf(`
		(function(require, module, exports) {
			%s
		})(require, module, module.exports);
	`, string(content))

	// 执行模块代码
	_, err = moduleVM.RunString(wrappedCode)
	if err != nil {
		return nil, err
	}

	// 缓存模块
	r.moduleCache[moduleName] = moduleObject.Get("exports").(*goja.Object)
	return moduleObject.Get("exports"), nil
}

func (r *Runtime) RunString(code string) (goja.Value, error) {
	return r.vm.RunString(code)
}
