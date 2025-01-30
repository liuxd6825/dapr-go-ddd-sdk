package server

import (
	"github.com/dop251/goja"
)

type Proxy struct {
	server *Server
	tpl    goja.Value
	vm     *goja.Runtime
}

func NewProxy(server *Server, vm *goja.Runtime) *Proxy {
	return &Proxy{
		server: server,
		vm:     vm,
		tpl:    vm.ToValue(server.tpl),
	}
}

// Get 方法：获取键对应的值
func (s *Proxy) Get(name string) goja.Value {
	var res = goja.Undefined()
	switch name {
	case "tpl":
		res = s.tpl
	case "workPath":
		res = s.vm.ToValue(s.server.FsOpts().WorkPath)
	default:
		if value, exists := s.server.RunValues().Get(name); exists {
			res = s.vm.ToValue(value)
		}
	}
	return res
}

// Set 方法：设置键值
func (s *Proxy) Set(name string, val goja.Value) bool {
	s.server.RunValues().Set(name, val.Export())
	return true
}

// Has 方法：检查键是否存在
func (s *Proxy) Has(name string) bool {
	exists := s.server.RunValues().Has(name)
	return exists
}

// Delete 方法：删除键
func (s *Proxy) Delete(name string) bool {
	s.server.RunValues().Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *Proxy) Keys() []string {
	return s.server.RunValues().Keys()
}

func (s *Proxy) InitVM(vm *goja.Runtime) error {
	return s.server.InitVM(vm)
}
