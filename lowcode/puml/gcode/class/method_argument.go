package class

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"

type MethodArgument struct {
	method   *Method
	Name     string `json:"name,omitempty"`
	DataType string `json:"dataType,omitempty"`
}

func NewMethodArgument(method *Method) *MethodArgument {
	return &MethodArgument{method: method}
}

func newMethodArgument(f *Method, str string) *MethodArgument {
	list := utils.StringSplit(str, " ")
	count := len(list)
	rVal := NewMethodArgument(f)
	if count == 1 {
		rVal.Name = utils.StringTrim(list[0])
		rVal.DataType = ""
	} else if count > 1 {
		rVal.Name = utils.StringTrim(list[0])
		rVal.DataType = utils.StringTrim(list[1])
	}
	return rVal
}
