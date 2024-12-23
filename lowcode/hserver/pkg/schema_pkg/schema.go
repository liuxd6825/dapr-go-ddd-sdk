package schema_pkg

import (
	"bytes"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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
	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := "schema.json"
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(s.server.GetSchemaLoader())

	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}

	s.cache.Set(fileUrl, sch)
	return sch
}
