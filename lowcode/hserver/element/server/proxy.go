package server

import (
	"github.com/dop251/goja"
)

type Proxy struct {
	server *Server
	vm     *goja.Runtime
	funcs  map[string]goja.Value
}

func NewProxy(server *Server, vm *goja.Runtime) *Proxy {
	p := &Proxy{
		vm:     vm,
		server: server,
	}
	p.init()
	return p
}

func (s *Proxy) init() {
	s.funcs = map[string]goja.Value{}
	s.funcs["init"] = s.vm.ToValue(s.server.Init)
	s.funcs["loadPkg"] = s.vm.ToValue(s.server.LoadPkg)
	s.funcs["workPath"] = s.vm.ToValue(s.server.FsOpts().WorkPath)
	s.funcs["logs"] = s.vm.ToValue(s.server.Logs)
	s.funcs["app"] = s.vm.ToValue(s.server.App)
}

// Get 方法：获取键对应的值
func (s *Proxy) Get(name string) goja.Value {
	var res = goja.Undefined()
	if v, ok := s.funcs[name]; ok {
		res = v
	} else if value, exists := s.server.RunValues().Get(name); exists {
		res = s.vm.ToValue(value)
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
