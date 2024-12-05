package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"strings"
)

func (r *Runtime) Require(fileName string) (val goja.Value, err error) {
	val, err = r.require(r.reader, fileName)
	return val, err
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

func (r *Runtime) getWorkPath() string {
	value := r.vm.Get("workPath")
	if value == nil {
		return ""
	}
	return value.String()
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
	workPath := r.getWorkPath()
	if strings.Contains(moduleName, "/definition/types/") {
		i := strings.LastIndex(moduleName, "/")
		moduleName = moduleName[i+1:]
		val = r.vm.Get(moduleName)
		if val == nil {
			return nil, fmt.Errorf("failed to load module %s", modulePath)
		}
		return val, nil
	}
	// 读取模块文件内容
	content, err := reader.ReadFile(modulePath, &fsopts.Options{WorkPath: workPath})
	if err != nil {
		return nil, fmt.Errorf("failed to load module %s: %v", modulePath, err)
	}

	content, err = TransformTSCodeToJS(string(content))
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
		if err, ok := err.(*goja.Exception); ok {
			return nil, fmt.Errorf("%s", err.String())
		}
		return nil, err
	}

	// 缓存模块
	pkg := NewPkgProxy(r.vm, NewPkgValue(moduleVM, exports))
	exportsObj := r.vm.NewDynamicObject(pkg)
	r.moduleCache[moduleName] = exportsObj
	return exportsObj, nil

	/*
		defVal := exports.ToObject(r.vm).Get("default")
		if defVal != nil {
			r.PrintValue(defVal)
		}
		pkg := NewPkgProxy(r.vm, NewPkgValue(moduleVM, defVal))

		exportsObj := r.vm.NewDynamicObject(pkg)
		r.moduleCache[moduleName] = exportsObj
		return exportsObj, nil
	*/
}

func (r *Runtime) PrintValue(val goja.Value) {
	eInst := val.ToObject(r.vm)

	// 获取自有属性
	ownKeys := eInst.Keys()
	fmt.Println("Own Properties:")
	for _, key := range ownKeys {
		fmt.Println("-", key)
	}

	// 获取原型方法
	proto := eInst.Prototype().ToObject(r.vm)
	protoKeys := proto.Keys()
	fmt.Println("Prototype Methods:")
	for _, key := range protoKeys {
		fmt.Println("-", key)
		method := proto.Get(key)
		fun, ok := goja.AssertFunction(method)
		if ok {
			fmt.Println("-", fun)
		}

	}

}
