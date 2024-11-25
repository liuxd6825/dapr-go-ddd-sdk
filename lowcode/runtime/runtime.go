package runtime

import (
	"github.com/dop251/goja"
	"reflect"
	"strings"
)

// Pool 封装 goja 的池
type Pool struct {
	userPool        bool
	vm              *goja.Runtime
	pool            chan *goja.Runtime // 使用 channel 实现池化，限制池大小
	fieldNameMapper goja.FieldNameMapper
}

// FieldNameMapper
type FieldNameMapper struct{}

var fieldNameMapper = &FieldNameMapper{}
var DefaultPoolSize = 5

// NewPool 创建一个带大小限制的 Runtime 池
func NewPool(userPool bool) *Pool {
	r := &Pool{
		userPool:        userPool,
		fieldNameMapper: &FieldNameMapper{},
	}
	if userPool {
		r.pool = make(chan *goja.Runtime, DefaultPoolSize)
		// 预热池：创建 poolSize 个 goja.Runtime 实例
		for i := 0; i < DefaultPoolSize; i++ {
			r.pool <- NewRuntime()
		}
	} else {
		r.pool = nil
		r.vm = NewRuntime()
	}
	return r
}

func NewRuntime() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(fieldNameMapper)
	return vm
}

// Run 执行 JavaScript 代码，并传递参数，返回结果
func (r *Pool) Run(code string, opts ...func(vm *goja.Runtime) error) (any, error) {
	var vm *goja.Runtime

	//  当不启用pool模式时
	if r.userPool {
		// 从池中获取一个实例（阻塞等待）
		select {
		case vm = <-r.pool:
			// 获取到实例
		default:
			// 创建新实例（仅在池耗尽时创建，避免阻塞过久）
			vm = NewRuntime()
		}

		// 确保实例归还池
		defer func() {
			select {
			case r.pool <- vm: // 归还实例
			default: // 池已满时丢弃实例
			}
		}()
	} else {
		vm = r.vm
	}

	return r.run(vm, code, opts...)
}

// run
//
//	@Description:
//	@receiver r
//	@param vm
//	@param code
//	@param opts
//	@return any
//	@return error
func (r *Pool) run(vm *goja.Runtime, code string, opts ...func(vm *goja.Runtime) error) (any, error) {
	err := r.initVM(vm, opts...)
	if err != nil {
		return nil, err
	}
	// 执行 JavaScript 代码
	result, err := vm.RunString(code)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result.Export(), nil
	}
	return result, err
}

// initVM
//
//	@Description:
//	@receiver r
//	@param vm
//	@param opts
//	@return error
func (r *Pool) initVM(vm *goja.Runtime, opts ...func(vm *goja.Runtime) error) error {
	for _, opt := range opts {
		if opt != nil {
			if err := opt(vm); err != nil {
				return err
			}
		}
	}
	return nil
}

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
