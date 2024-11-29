package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"regexp"
)

func RecoverError(e error, recover any) error {
	var err error
	if e != nil {
		err = e
	} else if recover != nil {
		if ve, ok := recover.(error); ok {
			err = ve
		} else if obj, ok := recover.(*goja.Object); ok {
			eObj := obj.Export()
			err = fmt.Errorf("unknown error %s", eObj)
		} else {
			err = fmt.Errorf("unknown error %s", recover)
		}
	}
	return err
}

// TransformTSCodeToJS
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformTSCodeToJS(tsCode string) (string, error) {
	// 正则表达式匹配 namespace 块
	//namespaceRegex := regexp.MustCompile(`(?s)namespace\s+\w+\s*{\s*(.*?)\s*}`)
	// 替换 namespace 块，仅保留其内部内容
	//tsCode = namespaceRegex.ReplaceAllString(tsCode, "$1")

	// 匹配以 *.d.ts 结尾的 import 语句
	dtsImportRegex := regexp.MustCompile(`(?m)^\s*import\s+.*?from\s+["'].*?\.d\.ts["'];?\s*$`)
	// 替换匹配的 import 语句为空字符串
	tsCode = dtsImportRegex.ReplaceAllString(tsCode, "")

	// 清理多余的空行
	emptyLinesRegex := regexp.MustCompile(`(?m)^\s*$\n`)
	tsCode = emptyLinesRegex.ReplaceAllString(tsCode, "")

	// 正则表达式匹配完整的 interface 定义并删除（包括尾部多余的换行和空格）
	interfaceRegex := regexp.MustCompile(`(?s)interface\s+\w+\s*{[^}]*}\s*`)
	// 正则表达式匹配 class 定义并删除（包括尾部多余的换行和空格）
	classRegex := regexp.MustCompile(`(?s)class\s+\w+(\s+extends\s+\w+)?\s*{[^}]*}\s*`)

	// 正则表达式匹配 let 变量定义中的类型注解，包括复杂类型
	letRegex := regexp.MustCompile(`\blet\s+(\w+)\s*:\s*[^=]+\s*=\s*`)

	// 正则表达式匹配类似 const params: Params; 的定义并删除
	constSpecificRegex := regexp.MustCompile(`(?m)^\s*const\s+\w+\s*:\s*\w+\s*;?\s*$`)

	// 正则表达式匹配所有变量定义（let），将其替换为 var
	//variableRegex := regexp.MustCompile(`\blet\b`)

	// 正则表达式匹配 TypeScript 类型标注并移除
	typeAnnotationRegex := regexp.MustCompile(`:\s*\w+`)

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

	// 清理多余的空行和缩进
	//cleanupRegex := regexp.MustCompile(`(?m)^\s*}\s*$`) // 匹配单独的 } 并删除
	//tsCode = cleanupRegex.ReplaceAllString(tsCode, "")

	// 删除多余空行
	emptyLinesRegex = regexp.MustCompile(`(?m)^\s*$\n`)
	tsCode = emptyLinesRegex.ReplaceAllString(tsCode, "")

	return tsCode, nil
}
