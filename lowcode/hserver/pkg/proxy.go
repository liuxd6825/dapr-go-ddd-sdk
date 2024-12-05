package pkg

import (
	"github.com/dop251/goja"
	"sort"
)

type Proxy struct {
	server   Server
	vm       *goja.Runtime
	values   map[string]any
	workPath string
}

func NewProxy(server Server, vm *goja.Runtime, workPath string, values map[string]any) *Proxy {
	if values == nil {
		values = make(map[string]any)
	}
	proxy := &Proxy{
		server:   server,
		vm:       vm,
		values:   values,
		workPath: workPath,
	}
	values["default"] = proxy
	return proxy
}

// Get 方法：获取键对应的值
func (s *Proxy) Get(name string) goja.Value {
	var res = goja.Undefined()
	if v, ok := s.values[name]; ok {
		return s.vm.ToValue(v)
	}
	return res
}

// Set 方法：设置键值
func (s *Proxy) Set(name string, val goja.Value) bool {
	s.values[name] = val
	return true
}

// Has 方法：检查键是否存在
func (s *Proxy) Has(name string) bool {
	_, ok := s.values[name]
	return ok
}

// Delete 方法：删除键
func (s *Proxy) Delete(name string) bool {
	delete(s.values, name)
	return true
}

// Keys 方法：获取所有键
func (s *Proxy) Keys() []string {
	// 获取键
	keys := make([]string, 0, len(s.values))
	for key := range s.values {
		keys = append(keys, key)
	}
	// 对键进行排序
	sort.Strings(keys)
	return keys
}

func (s *Proxy) InitVM(vm *goja.Runtime) error {
	return nil
}
