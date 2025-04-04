package runtime

import (
	"github.com/dop251/goja"
)

type PkgProxy struct {
	vm    *goja.Runtime
	value *PkgValue
}

func NewPkgProxy(vm *goja.Runtime, value *PkgValue) *PkgProxy {
	return &PkgProxy{
		vm:    vm,
		value: value,
	}
}

// Get 方法：获取键对应的值
func (s *PkgProxy) Get(name string) goja.Value {
	var res = s.value.Get(name)
	return res
}

// Set 方法：设置键值
func (s *PkgProxy) Set(name string, val goja.Value) bool {
	err := s.value.Set(name, val)
	if err != nil {
		return false
	}
	return true
}

// Has 方法：检查键是否存在
func (s *PkgProxy) Has(name string) bool {
	exists := s.value.Has(name)
	return exists
}

// Delete 方法：删除键
func (s *PkgProxy) Delete(name string) bool {
	return true
}

// Keys 方法：获取所有键
func (s *PkgProxy) Keys() []string {
	return s.value.Keys()
}
