package params_pkg

import (
	"github.com/dop251/goja"
	element2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/jsonschema/v6"
)

type ParamsPkg struct {
	server element2.Server
	cache  *types.CMap[*jsonschema.Schema]
}

func New(server element2.Server) *ParamsPkg {
	return NewParamsPkg(server)
}

func NewParamsPkg(server element2.Server) *ParamsPkg {
	return &ParamsPkg{
		server: server,
		cache:  types.NewCMap[*jsonschema.Schema](),
	}
}

func (s *ParamsPkg) NewProxy(vm *goja.Runtime, workPath string) *element2.Proxy {
	values := map[string]any{
		"loadFile": s.LoadFile,
	}
	return element2.NewProxy(s.server, vm, workPath, values)
}

func (s *ParamsPkg) LoadFile(fileUrl string, workPath string) *jsonschema.Schema {
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
	sch := schema.NewJsonSchemaWithBytes(fileUrl, data)
	s.cache.Set(fileUrl, sch)
	return sch
	/*
		var paramsType *common.ParamsType
		err := jsonutils.Unmarshal(data, &paramsType)
		if err != nil {
			panic(err)
		}
		s.cache.Set(fileUrl, paramsType)
		return paramsType
	*/
}
