package class

import (
	"testing"
)

func Test_ParseClassName(t *testing.T) {
	res := ParseClassName("class Edge << Use, Name , List>> extends Base ,Edit implements IBase,IEdit{", "class")
	t.Log(res)
}
