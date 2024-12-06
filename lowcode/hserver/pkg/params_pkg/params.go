package params_pkg

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/xtype"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
)

type ParamsPkg struct {
	server pkg.Server
	cache  *types.CMap[*xtype.ParamsType]
}

func New(server pkg.Server) *ParamsPkg {
	return NewParamsPkg(server)
}

func NewParamsPkg(server pkg.Server) *ParamsPkg {
	return &ParamsPkg{
		server: server,
		cache:  types.NewCMap[*xtype.ParamsType](),
	}
}

func (s *ParamsPkg) NewProxy(vm *goja.Runtime, workPath string) *pkg.Proxy {
	values := map[string]any{
		"loadFile": s.LoadFile,
	}
	return pkg.NewProxy(s.server, vm, workPath, values)
}

func (s *ParamsPkg) LoadFile(fileUrl string, workPath string) *xtype.ParamsType {
	if s.server.GetCacheEnable() {
		val, ok := s.cache.Get(fileUrl)
		if ok {
			return val
		}
	}
	data := s.server.GetFsm().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}
	var paramsType *xtype.ParamsType
	err := jsonutils.Unmarshal(data, &paramsType)
	if err != nil {
		panic(err)
	}
	s.cache.Set(fileUrl, paramsType)
	return paramsType
}
