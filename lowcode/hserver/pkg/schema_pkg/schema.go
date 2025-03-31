package schema_pkg

import (
	"bytes"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/spf13/afero"
)

var cache *types.CMap[*jsonschema.Schema] = types.NewCMap[*jsonschema.Schema]() //缓存
type SchemaPkg struct {
	server element.Server
}

func New(server element.Server) *SchemaPkg {
	return &SchemaPkg{
		server: server,
	}
}

func (s *SchemaPkg) openFile(fileUrl string, workPath string) ([]byte, error) {
	fs := s.server.SrcFs()
	data, err := afero.ReadFile(fs, fileUrl)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *SchemaPkg) LoadFile(fileUrl string, workPath string) *jsonschema.Schema {
	if s.server.CacheEnable() {
		val, ok := cache.Get(fileUrl)
		if ok {
			return val
		}
	}

	data := s.server.FsPkg().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}

	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := fileUrl
	compiler := schema.NewCompiler()
	compiler.UseLoader(s.server.SchemaLoader())

	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}
	/*
		if maps, ok := reader.(map[string]any); ok {
			if properties, ok := maps["properties"].(map[string]any); ok {
				for key, propVal := range properties {
					propVal.(map[string]interface{})["$schema"] = schemaFile
					types := propVal["type"]
					sch.Properties[key].Types = jsonschema.NewTypes()
				}
			}
		}*/

	cache.Set(fileUrl, sch)
	return sch
}
