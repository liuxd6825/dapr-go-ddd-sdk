package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Field struct {
	NameElement
	class  *Class
	Access Access `json:"access"`
}

func NewField(ctx context.Context, class *Class, line string, namespace string, notes []*Comment) *Field {
	field := &Field{
		class: class,
	}
	field.init(line)
	field.InitBase(field, line, "field", namespace, notes)

	return field
}

func (f *Field) init(line string) {
	/*
		{field} name string
	*/

	ok, str, newLine := utils.StringBetweenAll(line, "{", "}")
	if ok {
		list := strings.Split(str, " ")
		for _, s := range list {
			if s == "field" {
			}
		}
	}

	/*
		-	private 私有
		#	protected 受保护
		~	package private 包内可见
		+	public 公有
	*/
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

	list := strings.Split(utils.StringTrim(newLine), " ")
	count := len(list)
	if count >= 2 {
		f.Name = list[0]
		f.DataType = list[1]
	} else if count == 1 {
		f.Name = list[0]
		f.DataType = "string"
	}
}

func (f *Field) FUpperName() string {
	return utils.FirstUpperName(f.Name)
}

func (f *Field) FLowerName() string {
	return utils.FirstUpperName(f.Name)
}
