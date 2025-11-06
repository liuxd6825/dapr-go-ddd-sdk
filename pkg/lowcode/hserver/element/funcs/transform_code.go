package funcs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/runtime/transform"
)

var tsc *transform.Tsc

var babel *transform.Babel

// 正则表达式匹配第一个 function 的参数部分，仅修改第一个匿名函数
var funcRe = regexp.MustCompile(`\(function\s*\([^)]*\)`)

// 匹配以 *.d.ts 结尾的 import 语句
var dtsImportRegex = regexp.MustCompile(`(?m)^\s*import\s+.*?from\s+["'].*?\.d\.ts["'];?\s*$`)

// 清理多余的空行
var emptyLinesRegex = regexp.MustCompile(`(?m)^\s*$\n`)

// 正则表达式匹配完整的 interface 定义并删除（包括尾部多余的换行和空格）
var interfaceRegex = regexp.MustCompile(`(?s)interface\s+\w+\s*{[^}]*}\s*`)

// 正则表达式匹配 class 定义并删除（包括尾部多余的换行和空格）
var classRegex = regexp.MustCompile(`(?s)class\s+\w+(\s+extends\s+\w+)?\s*{[^}]*}\s*`)

// 正则表达式匹配 let 变量定义中的类型注解，包括复杂类型
var letRegex = regexp.MustCompile(`\blet\s+(\w+)\s*:\s*[^=]+\s*=\s*`)

// 正则表达式匹配类似 const params: Params; 的定义并删除
var constSpecificRegex = regexp.MustCompile(`(?m)^\s*const\s+\w+\s*:\s*\w+\s*;?\s*$`)

// 正则表达式匹配所有变量定义（let），将其替换为 var
var variableRegex = regexp.MustCompile(`\blet\b`)

// 正则表达式匹配 TypeScript 类型标注并移除
var typeAnnotationRegex = regexp.MustCompile(`:\s*\w+`)

// 定义正则表达式匹配 `//` 与 `@go-runtime` 之间有多个空格的情况，并删除下一行
var goRuntimeRegex = regexp.MustCompile(`(?m)^\s*//\s*@go-runtime.*\n.*\n`)

// TransformFromTypeScript
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformFromTypeScript(tsCode string, fileName string) ([]byte, []*FuncParam, error) {
	params := getParams(tsCode)

	// 替换为仅保留 function()
	tsCode = funcRe.ReplaceAllString(tsCode, `(function ()`)
	// 删除import语句
	tsCode = dtsImportRegex.ReplaceAllString(tsCode, ``)
	// 删除 // @go-runtime 行
	tsCode = goRuntimeRegex.ReplaceAllString(tsCode, "\n")
	if tsc == nil {
		tsc = transform.NewTsc()
	}
	es5Code, err := tsc.TransformEs5(tsCode, fileName)
	return es5Code, params, err
}

// TransformFromEs6
//
//	@Description:
//	@param es6Code
//	@param fileName
//	@return []byte
//	@return error
func TransformFromEs6(es6Code string, fileName string) ([]byte, []*FuncParam, error) {
	params := getParams(es6Code)
	// 替换为仅保留 function()
	es6Code = funcRe.ReplaceAllString(es6Code, `(function ()`)
	// 删除import语句
	es6Code = dtsImportRegex.ReplaceAllString(es6Code, ``)
	//
	es6Code = emptyLinesRegex.ReplaceAllString(es6Code, "")
	// 删除 // @go-runtime 行
	es6Code = goRuntimeRegex.ReplaceAllString(es6Code, "\n")
	// 删除 interface 定义
	es6Code = interfaceRegex.ReplaceAllString(es6Code, "")
	// 删除 class 定义
	es6Code = classRegex.ReplaceAllString(es6Code, "")
	// 删除特定的 const 定义
	es6Code = constSpecificRegex.ReplaceAllString(es6Code, "")
	// 替换为无类型注解的形式
	es6Code = letRegex.ReplaceAllString(es6Code, `let $1 = `)
	// 将 let 替换为 var
	es6Code = variableRegex.ReplaceAllString(es6Code, "var")
	// 移除类型标注
	es6Code = typeAnnotationRegex.ReplaceAllString(es6Code, "")
	//
	es6Code = emptyLinesRegex.ReplaceAllString(es6Code, "")
	// 替换匹配的内容为空
	es6Code = goRuntimeRegex.ReplaceAllString(es6Code, "\n")

	if babel == nil {
		var err error
		if babel, err = transform.NewBabel(); err != nil {
			return nil, nil, err
		}
	}
	codeBytes, resErr := babel.Transform(es6Code, fileName, false, nil)
	return codeBytes, params, resErr
}

// TransformCode
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformCode(tsCode string, fileName string, transType common.TransformType) ([]byte, []*FuncParam, error) {
	var codeBytes []byte
	var err error
	var params []*FuncParam
	switch transType {
	case common.TransformType_TypeScript:
		codeBytes, params, err = TransformFromTypeScript(tsCode, fileName)
	case common.TransformType_ES6:
		codeBytes, params, err = TransformFromEs6(tsCode, fileName)
	}
	return codeBytes, params, err
}

// 自定义 require 函数
func require(r *runtime.Runtime, reader fs.Reader, modulePath string) (val goja.Value, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	moduleName := modulePath
	// 检查缓存
	if cachedModule, ok := r.ModuleCache[modulePath]; ok {
		return cachedModule, nil
	}
	workPath := r.GetWorkPath()
	if strings.Contains(moduleName, "/definition/types/") {
		i := strings.LastIndex(moduleName, "/")
		moduleName = moduleName[i+1:]
		val = r.VM.Get(moduleName)
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

	pkgFileName := workPath + moduleName
	content, _, err = TransformCode(string(content), pkgFileName, common.TransformType_TypeScript)
	if err != nil {
		return nil, fmt.Errorf("failed to load module %s: %v", pkgFileName, err)
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
	_ = r.SetPkg(moduleVM, r.Pkg)
	// 绑定 require 函数，让模块内可以嵌套调用
	r.SetRequire(moduleVM, reader)

	// 包装模块代码，注入 require、module 和 exports
	code := fmt.Sprintf(`
	(function(require, module, exports) {
		%s
	})(require, module, module.exports);
	`, string(content))

	// 执行模块代码
	_, err = moduleVM.RunString(code)
	if err != nil {
		if err, ok := err.(*goja.Exception); ok {
			return nil, fmt.Errorf("%s", err.String())
		}
		return nil, err
	}

	// 缓存模块
	pkg := runtime.NewPkgProxy(r.VM, runtime.NewPkgValue(moduleVM, exports))
	exportsObj := r.VM.NewDynamicObject(pkg)
	r.ModuleCache[moduleName] = exportsObj
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
func getParams(code string) []*FuncParam {
	funIdx := strings.Index(code, "function")
	if funIdx == -1 {
		return []*FuncParam{}
	}
	code = code[funIdx:]
	// 查找括号部分，提取参数字符串
	startIndex := strings.Index(code, "(")
	endIndex := strings.Index(code, ")")
	resList := make([]*FuncParam, 0)
	if startIndex != -1 && endIndex != -1 {
		// 获取括号中的内容：name: string, age: number
		paramsString := code[startIndex+1 : endIndex]

		// 分割每个参数
		params := strings.Split(paramsString, ",")

		// 处理每个参数
		for _, param := range params {
			// 去除可能的空格
			param = strings.TrimSpace(param)
			// 查找参数的名称和类型
			colonIndex := strings.Index(param, ":")
			if colonIndex != -1 {
				paramName := strings.Trim(param[:colonIndex], " ")
				paramType := strings.Trim(param[colonIndex+1:], " ")
				resList = append(resList,
					&FuncParam{
						Name: paramName,
						Type: paramType,
					},
				)
			}
		}
	} else {
	}
	return resList
}
