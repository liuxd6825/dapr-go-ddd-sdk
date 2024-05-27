package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Method struct {
	NameElement
	class      *Class
	Access     Access            `json:"access,omitempty"`
	Args       []*MethodArgument `json:"args,omitempty"`
	Results    []*MethodResult   `json:"results,omitempty"`
	IsStatic   bool              `json:"isStatic"`
	IsAbstract bool              `json:"isAbstract"`
}

func NewMethod(ctx context.Context, class *Class, line string, notes []*Comment) *Method {
	method := &Method{
		class:  class,
		Access: AccessPublic,
	}
	method.InitBase(method, line, "method", class.GetNamespaceName(), notes)
	return method
}

func (f *Method) Init(line string) *Method {
	/*
		'' 保存
		{method} {static} {abstract} Save(cmd SaveCommand) error
		'' 按Id查询
		FindById(ctx context.Context, id string) (Person, error)
		'' 按Id查询
		FindPaging(ctx context.Context, filter string) ([]*Person, error)
		'' 按Id查询
		FindList(ctx context.Context, filter string) (persons []*Person, isFound bool, err error)
	*/

	/*
		{method} methods()
		{static} String id
		{abstract} void methods()
	*/

	ok, str, newLine := utils.StringBetweenAll(line, "{", "}")
	if ok {
		list := strings.Split(str, " ")
		for _, s := range list {
			if s == "method" {

			} else if s == "static" {
				f.IsStatic = true
			} else if s == "abstract" {
				f.IsAbstract = true
			}
		}
	}

	if strings.Contains(newLine, "-") {
		f.Access = AccessPrivate
		newLine = newLine[1:]
	} else if strings.Contains(newLine, "#") {
		f.Access = AccessProtected
		newLine = newLine[1:]
	} else if strings.Contains(newLine, "~") {
		f.Access = AccessPackagePrivate
		newLine = newLine[1:]
	} else if strings.Contains(newLine, "+") {
		f.Access = AccessPublic
		newLine = newLine[1:]
	} else {
		f.Access = AccessPublic
	}
	newLine = utils.StringTrim(newLine)
	name, residue := utils.SymbolStr(newLine, "(")
	if name != "" {
		f.Name = name
	} else {
		return f
	}

	// 处理传入参数
	argsStr, residue := utils.SymbolStr(residue, ")")
	argList := strings.Split(argsStr, ",")
	for _, argStr := range argList {
		arg := newMethodArgument(f, argStr)
		f.Args = append(f.Args, arg)
	}

	// 处理返回值（支持多个）
	if !strings.Contains(residue, "(") {
		rVal := newMethodResult(f, residue)
		f.Results = append(f.Results, rVal)
	} else {
		_, residue = utils.SymbolStr(residue, "(")
		resultStr, _ := utils.SymbolStr(residue, ")")
		resultList := strings.Split(resultStr, ",")
		for _, rItem := range resultList {
			rVal := newMethodResult(f, rItem)
			f.Results = append(f.Results, rVal)
		}
	}

	return f
}

func (f *Method) LenArgs() int {
	return len(f.Args)
}

func (f *Method) LenResults() int {
	return len(f.Results)
}

func (f *Method) HasArgs() bool {
	return len(f.Args) > 0
}

func (f *Method) HasResults() bool {
	return len(f.Results) > 0
}

func (f *Method) NotEndArg(idx int) bool {
	if idx < f.LenArgs()-1 {
		return true
	}
	return false
}

func (f *Method) NotEndResult(idx int) bool {
	if idx < f.LenResults()-1 {
		return true
	}
	return false
}

func (f *Method) IsPublic() bool {
	if f.Access == AccessPublic {
		return true
	}
	return false
}

func (f *Method) FUpperName() string {
	return utils.FirstUpperName(f.Name)
}

func (f *Method) FLowerName() string {
	return utils.FirstUpperName(f.Name)
}
