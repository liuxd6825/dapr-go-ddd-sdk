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
	var res = goja.Undefined()
	switch name {
	case "tpl":
		res = s.tpl
	case "workPath":
		res = s.vm.ToValue(s.server.fsOpts.WorkPath)
	default:
		if value, exists := s.server.runValues.Get(name); exists {
			res = s.vm.ToValue(value)
		}
	}
	return res
}

// Set 方法：设置键值
func (s *ServerProxy) Set(name string, val goja.Value) bool {
	s.server.runValues.Set(name, val.Export())
	return true
}

// Has 方法：检查键是否存在
func (s *ServerProxy) Has(name string) bool {
	exists := s.server.runValues.Has(name)
	return exists
}

// Delete 方法：删除键
func (s *ServerProxy) Delete(name string) bool {
	s.server.runValues.Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *ServerProxy) Keys() []string {
	return s.server.runValues.Keys()
}

func (s *ServerProxy) InitVM(vm *goja.Runtime) error {
	return s.server.InitVM(vm)
}
