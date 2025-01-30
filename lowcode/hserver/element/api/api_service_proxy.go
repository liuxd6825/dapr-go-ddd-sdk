package api

import "github.com/dop251/goja"

type ServiceProxy struct {
	service *ApiService
	vm      *goja.Runtime
}

func NewServiceProxy(service *ApiService, vm *goja.Runtime) *ServiceProxy {
	return &ServiceProxy{
		service: service,
		vm:      vm,
	}
}

func (s *ServiceProxy) InitVM(vm *goja.Runtime) error {
	return s.service.server.InitVM(vm)
}

// Get 方法：获取键对应的值
func (s *ServiceProxy) Get(name string) goja.Value {
	var res = goja.Undefined()
	switch name {
	case "workPath":
		res = s.vm.ToValue(s.service.FsOpts().WorkPath)
	default:
		if val, exists := s.service.data.Get(name); exists {
			if value, ok := val.(goja.Value); ok {
				res = value
			} else {
				res = s.vm.ToValue(value)
			}
		}
	}
	return res
}

// Set 方法：设置键值
func (s *ServiceProxy) Set(name string, val goja.Value) bool {
	s.service.data.Set(name, val)
	return true
}

// Has 方法：检查键是否存在
func (s *ServiceProxy) Has(name string) bool {
	exists := s.service.data.Has(name)
	return exists
}

// Delete 方法：删除键
func (s *ServiceProxy) Delete(name string) bool {
	s.service.data.Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *ServiceProxy) Keys() []string {
	return s.service.data.Keys()
}
