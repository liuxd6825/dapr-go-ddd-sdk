package hserver

import (
	"github.com/dop251/goja"
)

type ServerProxy struct {
	server *Server
	tpl    goja.Value
	vm     *goja.Runtime
}

func NewServerProxy(server *Server, vm *goja.Runtime) *ServerProxy {
	return &ServerProxy{
		server: server,
		vm:     vm,
		tpl:    vm.ToValue(server.tpl),
	}
}

// Get 方法：获取键对应的值
func (s *ServerProxy) Get(name string) goja.Value {
	if name == "tpl" {
		return s.tpl
	}
	if value, exists := s.server.data.Get(name); exists {
		return s.vm.ToValue(value)
	}
	return goja.Undefined()
}

// Set 方法：设置键值
func (s *ServerProxy) Set(name string, val goja.Value) bool {
	s.server.data.Set(name, val.Export())
	return true
}

// Has 方法：检查键是否存在
func (s *ServerProxy) Has(name string) bool {
	exists := s.server.data.Has(name)
	return exists
}

// Delete 方法：删除键
func (s *ServerProxy) Delete(name string) bool {
	s.server.data.Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *ServerProxy) Keys() []string {
	return s.server.data.Keys()
}
