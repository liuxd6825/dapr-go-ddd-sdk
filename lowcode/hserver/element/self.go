package element

import "github.com/dop251/goja"

type Self interface {
	SetSelfVMValue(name string, vm *goja.Runtime) error
}
