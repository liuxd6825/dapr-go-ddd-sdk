package schema_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/jsonschema/v6"
)

type SchemaPkg struct {
	server pkg.Server
	cache  *types.CMap[*jsonschema.Schema] //缓存
}

func New(server pkg.Server) *SchemaPkg {
	return &SchemaPkg{
		server: server,
		cache:  types.NewCMap[*jsonschema.Schema](),
	}
}

func (s *SchemaPkg) LoadFile(fileUrl string, workPath string) *jsonschema.Schema {
	if s.server.GetCacheEnable() {
		val, ok := s.cache.Get(fileUrl)
		if ok {
			return val
		}
	}
	data := s.server.GetFsPkg().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}
	var schema *jsonschema.Schema
	err := jsonutils.Unmarshal(data, &schema)
	if err != nil {
		panic(err)
	}

	s.cache.Set(fileUrl, schema)
	return schema
}
