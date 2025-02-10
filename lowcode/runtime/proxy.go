package runtime

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"reflect"
)

type Proxy[T any] struct {
	obj    T
	vm     *goja.Runtime
	values map[string]any
}

func NewProxy[T any](obj T, vm *goja.Runtime) *Proxy[T] {
	p := &Proxy[T]{
		obj:    obj,
		vm:     vm,
		values: make(map[string]any),
	}
	err := p.AddMethods(obj)
	if err != nil {
		panic(err)
	}
	return p
}

func (s *Proxy[T]) Values() map[string]any {
	return s.values
}

func (s *Proxy[T]) VM() *goja.Runtime {
	return s.vm
}

func (s *Proxy[T]) Obj() T {
	return s.obj
}

// AddMethods 将结构体的所有方法添加到goja.Runtime中
func (s *Proxy[T]) AddMethods(obj interface{}) error {
	// 获取结构体的反射类型
	t := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	// 获取结构体的方法数量
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fn := s.GetFunc(objV, &method)
		funcName := stringutils.FirstLower(method.Name)
		s.values[funcName] = fn
	}
	return nil
}

func (s *Proxy[T]) GetFunc(obj reflect.Value, method *reflect.Method) any {
	// 获取方法名
	methodName := method.Name
	// 将结构体方法转换为一个 Go 函数
	fn := func(args ...any) goja.Value {
		fmt.Println(methodName)
		var res []reflect.Value
		// 根据方法的签名调用对应的函数
		// 如果方法没有参数
		if method.Type.NumIn() == 1 {
			// 直接调用方法，无参数
			res = obj.MethodByName(methodName).Call(nil)
		} else if method.Type.NumIn() == 2 {
			// 如果方法有一个参数
			var argVals []reflect.Value
			for _, arg := range args {
				argVals = append(argVals, reflect.ValueOf(arg))
			}
			// 通过 reflect 调用方法
			res = obj.MethodByName(methodName).Call(argVals)
		}
		for _, item := range res {
			v := item.Interface()
			if err, ok := v.(error); ok {
				panic(err)
			}
			return s.vm.ToValue(v)
		}
		return goja.Undefined()
	}
	return fn
}
