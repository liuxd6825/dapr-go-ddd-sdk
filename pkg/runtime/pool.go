package runtime

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
)

// Pool 封装 goja 的池
type Pool struct {
	userPool bool
	vm       *Runtime
	pool     chan *Runtime // 使用 channel 实现池化，限制池大小
	reader   fs.Reader
	pkg      *types.CMap[any]
}

type RunOptions = func(vm *goja.Runtime) error

// NewPool 创建一个带大小限制的 Runtime 池
func NewPool(userPool bool, reader fs.Reader, pkg *types.CMap[any]) *Pool {
	r := &Pool{
		userPool: true,
		reader:   reader,
		pkg:      pkg,
	}
	if userPool {
		r.pool = make(chan *Runtime, DefaultPoolSize)
		// 预热池：创建 poolSize 个 goja.Runtime 实例
		for i := 0; i < DefaultPoolSize; i++ {
			r.pool <- NewRuntime(reader, pkg)
		}
	} else {
		r.pool = nil
		r.vm = NewRuntime(reader, pkg)
	}
	return r
}

// Run 执行 JavaScript 代码，并传递参数，返回结果
func (p *Pool) Run(code string, opts ...RunOptions) (resVal any, err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	var vm *Runtime

	//  当不启用pool模式时
	if p.userPool {
		// 从池中获取一个实例（阻塞等待）
		select {
		case vm = <-p.pool:
		// 获取到实例
		default:
			// 创建新实例（仅在池耗尽时创建，避免阻塞过久）
			vm = NewRuntime(p.reader, p.pkg)
		}

		// 确保实例归还池
		defer func() {
			select {
			case p.pool <- vm: // 归还实例
			default: // 池已满时丢弃实例
			}
		}()
	} else {
		vm = p.vm
	}
	return p.run(vm, code, opts...)
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
func (p *Pool) run(runtime *Runtime, code string, opts ...RunOptions) (dat any, err error) {

	defer func() {
		err = RecoverError(err, recover())
	}()
	for _, opt := range opts {
		if opt != nil {
			if err := opt(runtime.VM); err != nil {
				return nil, err
			}
		}
	}

	// 执行 JavaScript 代码
	result, err := runtime.RunString(code)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result.Export(), nil
	}
	return result, err
}
