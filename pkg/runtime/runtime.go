package runtime

import (
	"reflect"
	"strings"

	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
)

type Runtime struct {
	VM          *goja.Runtime
	ModuleCache map[string]*goja.Object
	reader      fs.Reader
	Pkg         *types.CMap[any]
}

// 模块缓存

var DefaultPoolSize = 5
var fieldNameMapper = &FieldNameMapper{}

func NewRuntime(reader fs.Reader, pkg *types.CMap[any]) *Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(fieldNameMapper)
	r := &Runtime{
		VM:          vm,
		ModuleCache: map[string]*goja.Object{},
		reader:      reader,
		Pkg:         pkg,
	}
	r.SetRequire(r.VM, reader)
	r.SetPkg(r.VM, r.Pkg)
	return r
}

func (r *Runtime) SetPkg(vm *goja.Runtime, pkg *types.CMap[any]) error {
	return vm.Set("pkg", vm.NewDynamicObject(NewObjectProxy(vm, pkg)))
}

func (r *Runtime) RunString(code string) (goja.Value, error) {
	return r.VM.RunString(code)
}

func (r *Runtime) GetVM() *goja.Runtime {
	return r.VM
}

func (r *Runtime) Set(name string, val any) error {
	return r.VM.Set(name, val)
}

func (r *Runtime) Get(name string) goja.Value {
	return r.VM.Get(name)
}

func (r *Runtime) NewObject() *goja.Object {
	return r.VM.NewObject()
}

func (r *Runtime) NewDynamicObject(val goja.DynamicObject) *goja.Object {
	return r.VM.NewDynamicObject(val)
}

func (r *Runtime) NewArray(item ...any) *goja.Object {
	return r.VM.NewArray(item...)
}

func (r *Runtime) NewDynamicArray(a goja.DynamicArray) *goja.Object {
	return r.VM.NewDynamicArray(a)
}

type FieldNameMapper struct{}

// FieldName 映射字段名称
func (m *FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return m.lowerFirstLetter(f.Name) // 将字段名转换为小写
}

// MethodName 映射方法名称
func (m *FieldNameMapper) MethodName(t reflect.Type, mtd reflect.Method) string {
	return m.lowerFirstLetter(mtd.Name) // 将方法名转换为小写
}

// 辅助函数：将首字母转换为小写
func (m *FieldNameMapper) lowerFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
