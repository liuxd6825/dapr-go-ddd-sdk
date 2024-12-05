package transform

import "regexp"

var tsc *Tsc

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
//variableRegex := regexp.MustCompile(`\blet\b`)

// 正则表达式匹配 TypeScript 类型标注并移除
var typeAnnotationRegex = regexp.MustCompile(`:\s*\w+`)

// Transform
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func Transform(tsCode string) ([]byte, error) {

	// 替换为仅保留 function()
	tsCode = funcRe.ReplaceAllString(tsCode, `(function ()`)
	/*
		// 替换匹配的 import 语句为空字符串
		tsCode = dtsImportRegex.ReplaceAllString(tsCode, "")
		//
		tsCode = emptyLinesRegex.ReplaceAllString(tsCode, "")
		// 删除 interface 定义
		tsCode = interfaceRegex.ReplaceAllString(tsCode, "")
		// 删除 class 定义
		tsCode = classRegex.ReplaceAllString(tsCode, "")
		// 删除特定的 const 定义
		tsCode = constSpecificRegex.ReplaceAllString(tsCode, "")
		// 替换为无类型注解的形式
		tsCode = letRegex.ReplaceAllString(tsCode, `let $1 = `)
		// 将 let 替换为 var
		//tsCode = variableRegex.ReplaceAllString(tsCode, "var")
		// 移除类型标注
		tsCode = typeAnnotationRegex.ReplaceAllString(tsCode, "")

		tsCode = emptyLinesRegex.ReplaceAllString(tsCode, "")
	*/
	if tsc == nil {
		tsc = NewTsc()
	}

	es5Code, err := tsc.TransformEs5(tsCode)
	return es5Code, err
}
