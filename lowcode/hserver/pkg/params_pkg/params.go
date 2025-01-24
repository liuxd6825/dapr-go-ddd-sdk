package params_pkg

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
)

type ParamsPkg struct {
	server element.Server
	cache  *types.CMap[*common.ParamsType]
}

func New(server element.Server) *ParamsPkg {
	return NewParamsPkg(server)
}

func NewParamsPkg(server element.Server) *ParamsPkg {
	return &ParamsPkg{
		server: server,
		cache:  types.NewCMap[*common.ParamsType](),
	}
}

func (s *ParamsPkg) NewProxy(vm *goja.Runtime, workPath string) *element.Proxy {
	values := map[string]any{
		"loadFile": s.LoadFile,
	}
	return element.NewProxy(s.server, vm, workPath, values)
}

func (s *ParamsPkg) LoadFile(fileUrl string, workPath string) *common.ParamsType {
	if s.server.CacheEnable() {
		val, ok := s.cache.Get(fileUrl)
		if ok {
			return val
		}
	}
	data := s.server.FsPkg().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}
	var paramsType *common.ParamsType
	err := jsonutils.Unmarshal(data, &paramsType)
	if err != nil {
		panic(err)
	}
	s.cache.Set(fileUrl, paramsType)
	return paramsType
}
