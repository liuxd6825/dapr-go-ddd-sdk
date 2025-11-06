package runtime

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
)

type ObjectProxy struct {
	vm     *goja.Runtime
	values *types.CMap[any]
}

func NewObjectProxy(vm *goja.Runtime, values *types.CMap[any]) *ObjectProxy {
	return &ObjectProxy{
		vm:     vm,
		values: values,
	}
}

// Get 方法：获取键对应的值
func (s *ObjectProxy) Get(name string) goja.Value {
	if val, ok := s.values.Get(name); ok {
		return s.vm.ToValue(val)
	}
	return goja.Undefined()
}

// Set 方法：设置键值
func (s *ObjectProxy) Set(name string, val goja.Value) bool {
	if val != nil && val != goja.Undefined() {
		s.values.Set(name, val.Export())
	} else {
		s.values.Set(name, goja.Undefined())
	}
	return true
}

// Has 方法：检查键是否存在
func (s *ObjectProxy) Has(name string) bool {
	exists := s.values.Has(name)
	return exists
}

// Delete 方法：删除键
func (s *ObjectProxy) Delete(name string) bool {
	s.values.Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *ObjectProxy) Keys() []string {
	return s.values.Keys()
}
