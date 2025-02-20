package transform

import (
	"regexp"
	"strings"
)

var tsc *Tsc

var babel *Babel

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
		tsc = NewTsc()
	}
	es5Code, err := tsc.TransformEs5(tsCode, fileName)
	return es5Code, params, err
}

type FuncParam struct {
	Name string
	Type string
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
				//fmt.Printf("Parameter: %s, Type: %s\n", paramName, paramType)
			}
		}
	} else {
		//fmt.Println("No parameters found.")
	}
	return resList
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
		if babel, err = NewBabel(); err != nil {
			return nil, nil, err
		}
	}
	codeBytes, resErr := babel.Transform(es6Code, fileName, false, nil)
	return codeBytes, params, resErr
}
